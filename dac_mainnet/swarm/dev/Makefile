.PHONY: build cluster test

default: build

build:
	go build -o bin/swarm    github.com/davinciproject/davinci_coin/dac_mainnet/cmd/swarm
	go build -o bin/geth     github.com/davinciproject/davinci_coin/dac_mainnet/cmd/geth
	go build -o bin/bootnode github.com/davinciproject/davinci_coin/dac_mainnet/cmd/bootnode

cluster: build
	scripts/boot-cluster.sh

test:
	go test -v github.com/davinciproject/davinci_coin/dac_mainnet/swarm/...
