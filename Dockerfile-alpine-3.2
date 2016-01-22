FROM alpine:3.2
MAINTAINER Will <zhguo.dataman-inc.com>

RUN mkdir -p /config

ADD config/haproxy_template.gateway.centos.cfg /config/haproxy_template.gateway.cfg
ADD config/haproxy_template.proxy.centos.cfg /config/haproxy_template.proxy.cfg
ADD config/production.proxy.json /config/production.proxy.json
ADD config/production.gateway.json /config/production.gateway.json
ADD config/production.example.json /config/production.example.json

ADD . /gopath/src/github.com/QubitProducts/bamboo
ADD haproxy /usr/share/haproxy
ADD builder/supervisord.conf /etc/supervisord.conf
ADD builder/run.sh /run.sh
ADD builder/buildBamboo.sh /buildBamboo.sh
WORKDIR /

RUN sh /buildBamboo.sh

VOLUME /var/log/supervisor
VOLUME /config

EXPOSE 80

CMD sh /run.sh   
