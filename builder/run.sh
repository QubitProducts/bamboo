#!/bin/bash
if [[ -n $BAMBOO_DOCKER_AUTO_HOST ]]; then
sed -i "s/^.*Endpoint\": \"\(http:\/\/haproxy-ip-address:8000\)\".*$/    \"EndPoint\": \"http:\/\/$HOST:8000\",/" \
    ${CONFIG_PATH:=config/production.example.json}
fi
haproxy -f /etc/haproxy/haproxy.cfg -p /var/run/haproxy.pid
/usr/bin/supervisord
