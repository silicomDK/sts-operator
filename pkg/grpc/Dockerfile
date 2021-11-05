ARG RUNTIME_IMAGE

FROM $RUNTIME_IMAGE as runtime
FROM golang:1.16 as builder

WORKDIR /workspace

COPY grpc/go.mod go.mod
COPY grpc/go.sum go.sum
COPY grpc/tsynctl tsynctl
RUN go mod download

COPY grpc/main.go main.go
COPY --from=runtime /lib64/libtsync.so /lib64/libtsync.so
COPY ./tsync_api/libs/libtsync/include/libtsync.h /usr/include/libtsync.h

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -o gprc-tsyncd main.go

COPY grpc/client.go client.go
RUN GOOS=linux GOARCH=amd64 go build -a -o grpc-client client.go

FROM runtime
COPY --from=builder /workspace/gprc-tsyncd /
COPY --from=builder /workspace/grpc-client /
ENTRYPOINT ["/gprc-tsyncd"]
