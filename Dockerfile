FROM ubuntu:14.04

RUN apt-get update -y
RUN apt-get install -y software-properties-common
RUN add-apt-repository ppa:vbernat/haproxy-1.5
RUN apt-get update -y
RUN apt-get install -y haproxy
RUN apt-get install -y golang
RUN apt-get install -y git
RUN apt-get install -y mercurial && rm -rf /var/lib/apt/lists/*

RUN apt-get update && apt-get install -y openssh-server supervisor
RUN mkdir -p /var/run/sshd /var/log/supervisor /root/.ssh && chmod 600 /root/.ssh
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
ADD authorized_keys /root/.ssh/authorized_keys

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

CMD /usr/bin/supervisord
