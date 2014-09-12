FROM ubuntu:14.04

RUN apt-get update -y
RUN apt-get install -y software-properties-common
RUN add-apt-repository ppa:vbernat/haproxy-1.5
RUN apt-get update -y
RUN apt-get install -y haproxy
RUN apt-get install -y golang
RUN apt-get install -y git
RUN apt-get install -y mercurial

ENV GOPATH /opt/go

RUN go get github.com/tools/godep
RUN go get -t github.com/smartystreets/goconvey

ADD . /opt/go/src/github.com/QubitProducts/bamboo
WORKDIR /opt/go/src/github.com/QubitProducts/bamboo
RUN /opt/go/bin/godep restore
RUN go build
RUN ln -s /opt/go/src/github.com/QubitProducts/bamboo /var/bamboo

RUN mkdir -p /run/haproxy

EXPOSE 8000
EXPOSE 80

CMD ["--help"]
ENTRYPOINT ["/var/bamboo/bamboo"]