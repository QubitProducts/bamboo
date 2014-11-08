FROM ubuntu:14.04

RUN apt-get update -y && apt-get install -y software-properties-common
RUN add-apt-repository ppa:vbernat/haproxy-1.5
RUN apt-get update -y && apt-get install -y haproxy golang git mercurial supervisor && rm -rf /var/lib/apt/lists/*

ENV GOPATH /opt/go

RUN go get github.com/tools/godep
RUN go get -t github.com/smartystreets/goconvey

ADD . /opt/go/src/github.com/QubitProducts/bamboo
WORKDIR /opt/go/src/github.com/QubitProducts/bamboo
RUN /opt/go/bin/godep restore
RUN go build
RUN ln -s /opt/go/src/github.com/QubitProducts/bamboo /var/bamboo

RUN mkdir -p /run/haproxy

ADD builder/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
ADD builder/run.sh /run.sh
RUN chmod +x /run.sh

RUN mkdir -p /var/log/supervisor
VOLUME /var/log/supervisor

EXPOSE 80 8000

CMD /run.sh
