FROM golang:1.6.2

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update -yqq && \
    apt-get install -yqq software-properties-common && \
    apt-get install -yqq haproxy git mercurial supervisor && \
    rm -rf /var/lib/apt/lists/*

ADD . /go/src/github.com/QubitProducts/bamboo
ADD builder/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
ADD builder/run.sh /run.sh

WORKDIR /go/src/github.com/QubitProducts/bamboo

RUN go get github.com/tools/godep && \
    go get -t github.com/smartystreets/goconvey && \
    go build && \
    ln -s /go/src/github.com/QubitProducts/bamboo /var/bamboo && \
    mkdir -p /run/haproxy && \
    mkdir -p /var/log/supervisor

VOLUME /var/log/supervisor

RUN apt-get clean && \
    rm -rf /tmp/* /var/tmp/* && \
    rm -rf /var/lib/apt/lists/* && \
    rm -f /etc/dpkg/dpkg.cfg.d/02apt-speedup && \
    rm -f /etc/ssh/ssh_host_*

EXPOSE 80 8000

CMD /run.sh
