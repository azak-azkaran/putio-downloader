VERSION := $(shell git describe --always --long --dirty)
all: install

fetch:
	@go get -u github.com/stretchr/testify
	@go get -u ./...

build: fetch
	@echo Building to current folder
	go build -i -v -ldflags="-X main.version=${VERSION}" 

docker: test
	docker build -t azakazkaran/putio-downloader .

install: build
	@echo Installing to ${GOPATH}/bin
	go install

test: fetch
	@echo Running tests
	go test 

coverage: test
	@echo Running Test with Coverage export
	go test -coverprofile=cover.out
	go test -json > report.json
	#go test github.com/azak-azkaran/cascade/utils -coverprofile=./utils/cover.out
	#go test github.com/azak-azkaran/cascade/utils -json > ./utils/report.json
	#cd ../

clean:
	go clean

push-image:
	docker push azakazkaran/putio-downloader:latest
