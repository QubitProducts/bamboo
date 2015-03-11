FROM ubuntu:14.04

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get install -yqq software-properties-common && \
    add-apt-repository -y ppa:vbernat/haproxy-1.5 && \
    apt-get update -yqq && \
    apt-get install -yqq haproxy golang git mercurial supervisor && \
    rm -rf /var/lib/apt/lists/*

ENV GOPATH /opt/go

RUN go get github.com/tools/godep && \
    go get -t github.com/smartystreets/goconvey

ADD . /opt/go/src/github.com/QubitProducts/bamboo
WORKDIR /opt/go/src/github.com/QubitProducts/bamboo
RUN /opt/go/bin/godep restore && \
    go build && \
    ln -s /opt/go/src/github.com/QubitProducts/bamboo /var/bamboo && \
    mkdir -p /run/haproxy

ADD builder/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
ADD builder/run.sh /run.sh
RUN chmod +x /run.sh

RUN mkdir -p /var/log/supervisor
VOLUME /var/log/supervisor

EXPOSE 80 8000

CMD /run.sh
