#!/usr/bin/env bash
set -ex

TAG=${TAG:-$TAG}

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
CODE_DIR=$(realpath $DIR/../..)

if [[ -z "$TAG" ]]; then
  echo "TAG shold be defined"
  exit 1
fi

VERSION_PLACEHOLDER="0.0.0"
sed -i '' "s/$VERSION_PLACEHOLDER/$TAG/g" $CODE_DIR/ddos-guard/main.go

# docker buildx build --push \
#   --platform linux/arm64/v8,linux/amd64 \
#   --tag yukels97/hello:latest \
#   -f $CODE_DIR/test/hello/Dockerfile $CODE_DIR

docker buildx build --push \
  --platform linux/arm64/v8,linux/amd64 \
  --tag yukels97/ddos-guard:latest \
  -f $CODE_DIR/ddos-guard/Dockerfile $CODE_DIR

docker buildx build --push \
  --platform linux/arm64/v8,linux/amd64 \
  --tag yukels97/ddos-guard:$TAG \
  -f $CODE_DIR/ddos-guard/Dockerfile $CODE_DIR