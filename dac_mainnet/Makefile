# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: gdac android ios gdac-cross swarm evm all test clean
.PHONY: gdac-linux gdac-linux-386 gdac-linux-amd64 gdac-linux-mips64 gdac-linux-mips64le
.PHONY: gdac-linux-arm gdac-linux-arm-5 gdac-linux-arm-6 gdac-linux-arm-7 gdac-linux-arm64
.PHONY: gdac-darwin gdac-darwin-386 gdac-darwin-amd64
.PHONY: gdac-windows gdac-windows-386 gdac-windows-amd64

GOBIN = $(shell pwd)/build/bin
GO ?= latest

gdac:
	build/env.sh go run build/ci.go install ./cmd/gdac
	@echo "Done building."
	@echo "Run \"$(GOBIN)/gdac\" to launch gdac."

all:
	build/env.sh go run build/ci.go install

android:
	build/env.sh go run build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/gdac.aar\" to use the library."

ios:
	build/env.sh go run build/ci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/Geth.framework\" to use the library."

test: all
	build/env.sh go run build/ci.go test

lint: ## Run linters.
	build/env.sh go run build/ci.go lint

clean:
	./build/clean_go_build_cache.sh
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go get -u golang.org/x/tools/cmd/stringer
	env GOBIN= go get -u github.com/kevinburke/go-bindata/go-bindata
	env GOBIN= go get -u github.com/fjl/gencodec
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go install ./cmd/abigen
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

# Cross Compilation Targets (xgo)

gdac-cross: gdac-linux gdac-darwin gdac-windows gdac-android gdac-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/gdac-*

gdac-linux: gdac-linux-386 gdac-linux-amd64 gdac-linux-arm gdac-linux-mips64 gdac-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/gdac-linux-*

gdac-linux-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/gdac
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/gdac-linux-* | grep 386

gdac-linux-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/gdac
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gdac-linux-* | grep amd64

gdac-linux-arm: gdac-linux-arm-5 gdac-linux-arm-6 gdac-linux-arm-7 gdac-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/gdac-linux-* | grep arm

gdac-linux-arm-5:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/gdac
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/gdac-linux-* | grep arm-5

gdac-linux-arm-6:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/gdac
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/gdac-linux-* | grep arm-6

gdac-linux-arm-7:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/gdac
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/gdac-linux-* | grep arm-7

gdac-linux-arm64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/gdac
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/gdac-linux-* | grep arm64

gdac-linux-mips:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/gdac
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/gdac-linux-* | grep mips

gdac-linux-mipsle:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/gdac
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/gdac-linux-* | grep mipsle

gdac-linux-mips64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/gdac
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/gdac-linux-* | grep mips64

gdac-linux-mips64le:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/gdac
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/gdac-linux-* | grep mips64le

gdac-darwin: gdac-darwin-386 gdac-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/gdac-darwin-*

gdac-darwin-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/gdac
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/gdac-darwin-* | grep 386

gdac-darwin-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/gdac
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gdac-darwin-* | grep amd64

gdac-windows: gdac-windows-386 gdac-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/gdac-windows-*

gdac-windows-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/gdac
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/gdac-windows-* | grep 386

gdac-windows-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/gdac
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gdac-windows-* | grep amd64
