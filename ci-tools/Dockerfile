# Base image: https://hub.docker.com/_/golang/
FROM golang:1.12-buster
MAINTAINER Sbit.io <soporte@sbit.io>

# Install golint
ENV GOPATH /go
ENV PATH ${GOPATH}/bin:$PATH

# Install golint
RUN go get -u golang.org/x/lint/golint

# Add LLVM's clang
RUN wget -O - https://apt.llvm.org/llvm-snapshot.gpg.key | apt-key add - \
 && echo "deb http://apt.llvm.org/buster/ llvm-toolchain-buster main" > /etc/apt/sources.list.d/llvm.list \
 && apt-get update -qq \
 && DEBIAN_FRONTEND=noninteractive \
    apt-get install -yqq -o=Dpkg::Use-Pty=0 \
      clang \
 && rm -rf /var/lib/apt/lists/* \
 && apt-get clean -yqq
