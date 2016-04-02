all: install

checks:
	@echo "Checking deps:"
	@(env bash $(PWD)/buildscripts/checkgopath.sh)

getdeps: checks
	@go get -u golang.org/x/tools/cmd/vet && echo "Installed vet:"
	@go get -u github.com/jteeuwen/go-bindata/... && echo "Installed go-bindata:"
	@go get -u github.com/elazarl/go-bindata-assetfs/... && echo "Installed go-bindata-assetfs:"

verifiers: vet fmt

vet:
	@echo "Running $@:"
	@GO15VENDOREXPERIMENT=1 go tool vet -all ./media-player
	@GO15VENDOREXPERIMENT=1 go tool vet -shadow=true ./media-player

fmt:
	@echo "Running $@:"
	@GO15VENDOREXPERIMENT=1 gofmt -s -l media-player

webassets:
	@echo "Generating $@:"
	@GO15VENDOREXPERIMENT=1 ${GOPATH}/bin/go-bindata-assetfs web/...
	@mv bindata_assetfs.go media-player/web-assets.go
	@gofmt -s -w -l media-player/web-assets.go
	@echo "Please commit media-player/web-assets.go"

build: getdeps verifiers

gomake-all: build
	@echo "Installing media-player:"
	@go install github.com/minio/minio-go-media-player/media-player

install: gomake-all

clean:
	@echo "Cleaning up all the generated files:"
	@rm -fv cover.out
	@find . -name '*.test' | xargs rm -fv
