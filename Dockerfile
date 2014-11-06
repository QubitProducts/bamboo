FROM ubuntu:14.04

RUN apt-get update -y && apt-get install -y software-properties-common
RUN add-apt-repository ppa:vbernat/haproxy-1.5
RUN apt-get update -y && apt-get install -y haproxy golang git mercurial supervisor openssh-server && rm -rf /var/lib/apt/lists/*

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
ADD builder/start-haproxy.sh /start-haproxy.sh
RUN chmod +x /start-haproxy.sh

RUN mkdir -p /var/run/sshd /var/log/supervisor /root/.ssh && chmod 600 /root/.ssh
ADD authorized_keys /root/.ssh/authorized_keys

EXPOSE 8000
EXPOSE 80

CMD /usr/bin/supervisord
