#!/bin/bash

set -ex
export VERSION=0.0.2
export IMAGE=quay.io/silicom/sts-operator:$VERSION
export IMG=$IMAGE
export BUNDLE_IMG=quay.io/silicom/sts-operator-bundle:$VERSION
export NAMESPACE=sts-silicom

rm -rf bundle
#podman build -t $IMG:latest .
make docker-build
make docker-push

make bundle
make bundle-build
make bundle-push
