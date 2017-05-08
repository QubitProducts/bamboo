FROM golang:1.8

ENV DEBIAN_FRONTEND noninteractive

RUN echo deb http://httpredir.debian.org/debian jessie-backports main | \
      sed 's/\(.*-backports\) \(.*\)/&@\1-sloppy \2/' | tr @ '\n' | \
      tee /etc/apt/sources.list.d/backports.list && \
    curl https://haproxy.debian.net/bernat.debian.org.gpg | \
      apt-key add - && \
    echo deb http://haproxy.debian.net jessie-backports-1.5 main | \
      tee /etc/apt/sources.list.d/haproxy.list

RUN apt-get update -yqq && \
    apt-get install -yqq software-properties-common && \
    apt-get install -yqq git mercurial supervisor && \
    apt-get install -yqq haproxy -t jessie-backports-1.5 && \
    rm -rf /var/lib/apt/lists/*

ADD builder/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
ADD builder/run.sh /run.sh

WORKDIR /go/src/github.com/QubitProducts/bamboo

RUN go get github.com/tools/godep && \
    go get -t github.com/smartystreets/goconvey

ADD . /go/src/github.com/QubitProducts/bamboo

RUN go build && \
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

