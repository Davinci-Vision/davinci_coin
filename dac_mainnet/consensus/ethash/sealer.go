// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package ethash

import (
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/davinciproject/davinci_coin/dac_mainnet/common"
	"github.com/davinciproject/davinci_coin/dac_mainnet/consensus"
	"github.com/davinciproject/davinci_coin/dac_mainnet/core/types"
	"github.com/davinciproject/davinci_coin/dac_mainnet/log"
	"github.com/davinciproject/davinci_coin/dac_mainnet/crypto"
)

// Seal implements consensus.Engine, attempting to find a nonce that satisfies
// the block's difficulty requirements.
func (ethash *Ethash) Seal(chain consensus.ChainReader, block *types.Block, stakeAmount *big.Int, stop <-chan struct{}) (*types.Block, error) {
	// If we're running a fake PoW, simply return a 0 nonce immediately
	if ethash.config.PowMode == ModeFake || ethash.config.PowMode == ModeFullFake {
		header := block.Header()
		header.Nonce, header.MixDigest = types.BlockNonce{}, common.Hash{}
		return block.WithSeal(header), nil
	}
	// If we're running a shared PoW, delegate sealing to it
	if ethash.shared != nil {
		return ethash.shared.Seal(chain, block, stakeAmount, stop)
	}
	// Create a runner and the multiple search threads it directs
	abort := make(chan struct{})
	found := make(chan *types.Block)

	ethash.lock.Lock()
	threads := ethash.threads
	if ethash.rand == nil {
		seed, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
		if err != nil {
			ethash.lock.Unlock()
			return nil, err
		}
		ethash.rand = rand.New(rand.NewSource(seed.Int64()))
	}
	ethash.lock.Unlock()
	if threads == 0 {
		threads = 1
	}
	if threads < 0 {
		threads = 0 // Allows disabling local mining without extra logic around local/remote
	}
	var pend sync.WaitGroup
	for i := 0; i < threads; i++ {
		pend.Add(1)
		if !chain.Config().IsByzantium(block.Number()) {
			go func(id int, nonce uint64) {
				defer pend.Done()
				ethash.mine(block, id, nonce, abort, found)
			}(i, uint64(ethash.rand.Int63()))
		} else {
			go func(id int) {
				defer pend.Done()
				ethash.minePOS(block, id, stakeAmount, chain.CurrentHeader().Time, abort, found)
			}(i)
		}
	}
	// Wait until sealing is terminated or a nonce is found
	var result *types.Block
	select {
	case <-stop:
		// Outside abort, stop all miner threads
		close(abort)
	case result = <-found:
		// One of the threads found a block, abort all others
		close(abort)
	case <-ethash.update:
		// Thread count was changed on user request, restart
		close(abort)
		pend.Wait()
		return ethash.Seal(chain, block, stakeAmount, stop)
	}
	// Wait for all miners to terminate and return the block
	pend.Wait()
	return result, nil
}

// mine is the actual proof-of-work miner that searches for a nonce starting from
// seed that results in correct final block difficulty.
func (ethash *Ethash) mine(block *types.Block, id int, seed uint64, abort chan struct{}, found chan *types.Block) {
	// Extract some data from the header
	var (
		header  = block.Header()
		hash    = header.HashNoNonce().Bytes()
		target  = new(big.Int).Div(maxUint256, header.Difficulty)
		number  = header.Number.Uint64()
		dataset = ethash.dataset(number)
	)
	// Start generating random nonces until we abort or find a good one
	var (
		attempts = int64(0)
		nonce    = seed
	)
	logger := log.New("miner", id)
	logger.Trace("Started ethash search for new nonces", "seed", seed)
search:
	for {
		select {
		case <-abort:
			// Mining terminated, update stats and abort
			logger.Trace("Ethash nonce search aborted", "attempts", nonce-seed)
			ethash.hashrate.Mark(attempts)
			break search

		default:
			// We don't have to update hash rate on every nonce, so update after after 2^X nonces
			attempts++
			if (attempts % (1 << 15)) == 0 {
				ethash.hashrate.Mark(attempts)
				attempts = 0
			}
			// Compute the PoW value of this nonce
			digest, result := hashimotoFull(dataset.dataset, hash, nonce)
			if new(big.Int).SetBytes(result).Cmp(target) <= 0 {
				// Correct nonce found, create a new header with it
				header = types.CopyHeader(header)
				header.Nonce = types.EncodeNonce(nonce)
				header.MixDigest = common.BytesToHash(digest)

				// Seal and return a block (if still needed)
				select {
				case found <- block.WithSeal(header):
					logger.Trace("Ethash nonce found and reported", "attempts", nonce-seed, "nonce", nonce)
				case <-abort:
					logger.Trace("Ethash nonce found but discarded", "attempts", nonce-seed, "nonce", nonce)
				}
				break search
			}
			nonce++
		}
	}
	// Datasets are unmapped in a finalizer. Ensure that the dataset stays live
	// during sealing so it's not unmapped while being read.
	runtime.KeepAlive(dataset)
}

// mine with POS algorithm
func (ethash *Ethash) minePOS(block *types.Block, id int, stakeAmount *big.Int, parentTime *big.Int, abort chan struct{}, found chan *types.Block) {
	
	// Extract some data from the header
	var (
		header  = block.Header()
	)
	logger := log.New("miner", id)
	logger.Trace("Started POS sealing")
search:
	for {
		select {
		case <-abort:
			// Mining terminated, update stats and abort
			logger.Trace("POS sealing aborted")
			break search

		default:
			header = types.CopyHeader(header)
			header.StakeAmount = stakeAmount

			difficulty := header.Difficulty
			coinbase := header.Coinbase
			parentHash := header.ParentHash

			if stakeAmount.Cmp(big.NewInt(0)) <= 0 {
				logger.Trace("mining with zero stake. mining suspended", "address", coinbase, "stakeAmount", stakeAmount)
				break search
			} else {
				result := new(big.Int).SetBytes(crypto.Keccak256(parentHash.Bytes(), coinbase.Bytes()))
				blockTime := new(big.Int).Div( new(big.Int).Mul(result,difficulty) , new(big.Int).Mul(stakeAmount,maxUint256) )
				blockTime = new(big.Int).Add(blockTime, big.NewInt(1))
				if blockTime.Cmp(big.NewInt(86400)) >= 0 {
					logger.Trace("will take more than a day to mine a block", "blockTime", blockTime)
					break search
				}

				newTime := new(big.Int).Add(parentTime, blockTime)
				intTime := newTime.Int64()

				now := time.Now().Unix()
				if intTime > now {
					time.Sleep(time.Duration(intTime - now) * time.Second)
				}
				header.Time = newTime

				// Seal and return a block (if still needed)
				select {
				case found <- block.WithSeal(header):
					logger.Trace("POS sealing succeeded and reported")
				case <-abort:
					logger.Trace("POS sealing succeeded but discarded")
				}
				break search
			}
		}
	}
}