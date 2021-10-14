# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2020 Intel Corporation

FROM golang

WORKDIR /workspace
COPY go.mod go.sum ./
RUN go mod download

COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o sts-plugin main.go

FROM registry.access.redhat.com/ubi8/ubi-minimal
RUN microdnf install pciutils iproute

COPY --from=0 /workspace/sts-plugin /

ENTRYPOINT ["/sts-plugin"]
