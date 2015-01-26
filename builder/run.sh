#!/bin/bash
if [[ -n $AUTO_BAMBOO_HOST ]]; then
sed -i "s/^.*Endpoint\": \"\(http:\/\/haproxy-ip-address:8000\)\".*$/    \"EndPoint\": \"$HOST\",/" \
    ${CONFIG_PATH:=config/production.example.json}
fi
haproxy -f /etc/haproxy/haproxy.cfg -p /var/run/haproxy.pid
/usr/bin/supervisord
