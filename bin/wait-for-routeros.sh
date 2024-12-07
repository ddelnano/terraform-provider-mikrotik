#!/bin/sh

routeros_host=${1:-127.0.0.1}
routeros_port=8728

echo "waiting for RouterOS (${routeros_host}:${routeros_port}) to be up and running"

for i in $(seq 1 60); do
    if nc -w1 ${routeros_host} ${routeros_port}; then
        exit 0;
    else
        printf "."
        sleep 1
    fi done;

exit 1
