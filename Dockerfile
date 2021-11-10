# Build the manager binary
FROM golang:1.17 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
COPY api/ api/
COPY grpc/ grpc/

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
