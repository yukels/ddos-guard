FROM ubuntu:18.04

RUN apt-get update \
    && apt-get install -y sudo \
    && apt-get install -y git \
    && apt-get install -y build-essential \
    && apt-get install -y jq \
    && apt-get install -y curl \
    && apt-get install -y wget \
    && apt-get install -y graphviz \
    && apt-get install -y gettext-base 

# Install GoLang
ENV GO_VERSION="1.20.1"
RUN buildOs="$(uname --kernel-name | tr '[:upper:]' '[:lower:]')"; \
    buildArch="$(uname --processor)"; \
    case "$buildArch" in \
        aarch64) buildArch='arm64';; \
        x86_64) buildArch='amd64';; \
    esac; \
    PLATFORM="${buildOs}-${buildArch}"; \
    wget https://golang.org/dl/go${GO_VERSION}.${PLATFORM}.tar.gz \
    && tar -xvf go${GO_VERSION}.${PLATFORM}.tar.gz \
    && mv go /usr/local \
    && rm go${GO_VERSION}.${PLATFORM}.tar.gz

ENV GOROOT="/usr/local/go"
ENV GOPATH="/workspace/go"
ENV PATH="${PATH}:${GOPATH}/bin:${GOROOT}/bin"    
