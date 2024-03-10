VERSION = "0.0.1"
BUILD_DATE = $(shell date | cat )
COMMIT = $(shell git log --pretty=format:"%H" | head -n 1)

compile-keeper:
	- rm -r keeper
	mkdir keeper
	cp ./keeper_config.example.json keeper/keeper_config.json
	GOOS=windows GOARCH=amd64 go build -o keeper/keeper.exe -ldflags "-X main.Commit=$(COMMIT) -X 'main.BuildTime=$(BUILD_DATE)' -X main.Version=$(VERSION)" cmd/keeper/main.go
	zip -r ./bin/keeper/keeper_win64.zip keeper
	rm keeper/keeper.exe
	GOOS=darwin GOARCH=amd64 go build -o keeper/keeper -ldflags "-X main.Commit=$(COMMIT) -X 'main.BuildTime=$(BUILD_DATE)' -X main.Version=$(VERSION)" cmd/keeper/main.go
	zip -r ./bin/keeper/keeper_darwin-amd64.zip keeper
	rm keeper/keeper
	GOOS=darwin GOARCH=arm64 go build -o keeper/keeper -ldflags "-X main.Commit=$(COMMIT) -X 'main.BuildTime=$(BUILD_DATE)' -X main.Version=$(VERSION)" cmd/keeper/main.go
	zip -r ./bin/keeper/keeper_darwin-arm.zip keeper
	rm keeper/keeper
	GOOS=linux GOARCH=amd64 go build -o keeper/keeper -ldflags "-X main.Commit=$(COMMIT) -X 'main.BuildTime=$(BUILD_DATE)' -X main.Version=$(VERSION)" cmd/keeper/main.go
	zip -r ./bin/keeper/keeper_linux.zip keeper
	rm -r keeper