VERSION = "0.0.1"
BUILD_DATE = $(shell date | cat )
COMMIT = $(shell git log --pretty=format:"%H" | head -n 1)

compile-keeper:
	- rm -r keeperbuild
	mkdir keeperbuild
	cp ./keeper_config.example.json keeperbuild/keeper_config.json
	GOOS=windows GOARCH=amd64 go build -o keeperbuild/keeper.exe -ldflags "-X main.Commit=$(COMMIT) -X 'main.BuildTime=$(BUILD_DATE)' -X main.Version=$(VERSION)" cmd/service/main.go
	zip -r ./bin/keeper/keeper_win64.zip keeperbuild
	rm keeperbuild/keeper.exe
	GOOS=darwin GOARCH=amd64 go build -o keeperbuild/keeper -ldflags "-X main.Commit=$(COMMIT) -X 'main.BuildTime=$(BUILD_DATE)' -X main.Version=$(VERSION)" cmd/service/main.go
	zip -r ./bin/keeper/keeper_darwin-amd64.zip keeperbuild
	rm keeperbuild/keeper
	GOOS=darwin GOARCH=arm64 go build -o keeperbuild/keeper -ldflags "-X main.Commit=$(COMMIT) -X 'main.BuildTime=$(BUILD_DATE)' -X main.Version=$(VERSION)" cmd/service/main.go
	zip -r ./bin/keeper/keeper_darwin-arm.zip keeperbuild
	rm keeperbuild/keeper
	GOOS=linux GOARCH=amd64 go build -o keeperbuild/keeper -ldflags "-X main.Commit=$(COMMIT) -X 'main.BuildTime=$(BUILD_DATE)' -X main.Version=$(VERSION)" cmd/service/main.go
	zip -r ./bin/keeper/keeper_linux.zip keeperbuild