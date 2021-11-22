# Build the manager binary
FROM golang:1.17 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
COPY api/ api/
COPY grpc/ grpc/

ADD	https://github.com/protocolbuffers/protobuf/releases/download/v3.15.8/protoc-3.15.8-linux-x86_64.zip protoc.zip
RUN apt update && \
    apt install unzip && \
    unzip protoc.zip && \
    mv bin/protoc /usr/bin/

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
RUN	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

RUN protoc --go_out=. \
           --go_opt=paths=source_relative \
           --go-grpc_out=. \
           --go-grpc_opt=paths=source_relative \
            grpc/tsynctl/grpc_tsynctl.proto


#RUN go mod tidy

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY controllers/ controllers/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager main.go

FROM registry.access.redhat.com/ubi8/ubi-minimal
RUN microdnf install pciutils kmod procps iproute nc
WORKDIR /
COPY --from=builder /workspace/manager .
COPY assets assets

ENTRYPOINT ["/manager"]
