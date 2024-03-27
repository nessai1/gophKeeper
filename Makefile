VERSION = "0.0.1"
BUILD_DATE = $(shell date | cat )
COMMIT = $(shell git log --pretty=format:"%H" | head -n 1)

compile-keeper:
	- rm -r keeper
	mkdir keeper
	cp ./keeper_config.example.json keeper/keeper_config.json
	$(MAKE) build-bin GOOS=windows GOARCH=amd64 BIN_NAME=keeper.exe BUILD_NAME=win64
	$(MAKE) build-bin GOOS=darwin GOARCH=amd64 BIN_NAME=keeper BUILD_NAME=darwin-amd64
	$(MAKE) build-bin GOOS=darwin GOARCH=arm64 BIN_NAME=keeper BUILD_NAME=darwin-arm
	$(MAKE) build-bin GOOS=linux GOARCH=amd64 BIN_NAME=keeper BUILD_NAME=linux
	rm -r keeper

build-bin:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o keeper/$(BIN_NAME) -ldflags "-X main.Commit=$(COMMIT) -X 'main.BuildTime=$(BUILD_DATE)' -X main.Version=$(VERSION)" cmd/keeper/main.go
	zip -r ./bin/keeper/keeper_$(BUILD_NAME).zip keeper
	rm keeper/$(BIN_NAME)