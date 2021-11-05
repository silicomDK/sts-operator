IMAGE_REGISTRY ?= quay.io/silicom
BUILDTOOL ?= docker
SILTSYNC_VER = $(shell cat ../VERSION)

#PATH="$$PATH:$$(go env GOPATH)/bin"

all: docker-build

protobuf:
	- mkdir downloads
	curl -sL https://github.com/protocolbuffers/protobuf/releases/download/v3.15.8/protoc-3.15.8-linux-x86_64.zip -o downloads/protoc.zip
	cd downloads && unzip protoc.zip
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

lib:
	/bin/bash lib.sh

docker-build: protoc
	cd ../ && $(BUILDTOOL) build . --build-arg=RUNTIME_IMAGE=$(IMAGE_REGISTRY)/siltsync:$(SILTSYNC_VER) -f grpc/Dockerfile -t $(IMAGE_REGISTRY)/grpc-tsyncd:$(SILTSYNC_VER)

docker-push:
	$(BUILDTOOL) push $(IMAGE_REGISTRY)/grpc-tsyncd:$(SILTSYNC_VER)

protoc:
	export PATH=$$(go env GOPATH)/bin:$$PATH; \
	downloads/bin/protoc --go_out=. --go_opt=paths=source_relative \
    					 --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    					 tsynctl/tsynctl.proto
clean:
	rm -rf downloads