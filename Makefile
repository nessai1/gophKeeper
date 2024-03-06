BUILD_DATE = $(shell date | cat )
COMMIT = $(shell git log --pretty=format:"%H" | head -n 1)

compile-service:
	- rm -r servicebuild
	mkdir servicebuild
	cp -r ./migrations servicebuild
	go build -o servicebuild/service -ldflags "-X main.Commit=$(COMMIT) -X 'main.BuildTime=$(BUILD_DATE)'" cmd/service/main.go
	zip -r ./bin/server/x64.zip servicebuild
	rm -r servicebuild

compile-keeper:
	echo $(VL)