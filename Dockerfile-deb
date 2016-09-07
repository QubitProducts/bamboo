FROM golang:1.7.0

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update -yqq
RUN apt-get install -yqq software-properties-common curl git mercurial ruby-dev gcc make rubygems
RUN gem install fpm

ADD . /go/src/github.com/QubitProducts/bamboo

WORKDIR /go/src/github.com/QubitProducts/bamboo
RUN go get github.com/tools/godep && \
    go get -t github.com/smartystreets/goconvey && \
    go build

RUN mkdir -p /output
VOLUME /output

ENV USER root
CMD builder/build.sh && mv output/bamboo*.deb /output
