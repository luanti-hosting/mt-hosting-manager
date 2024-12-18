#!/bin/bash
set -e
export DEBIAN_FRONTEND=noninteractive

cd `dirname $0`

test -f "APT_STAGE1" ||{
    apt-get update
    sleep 5
    apt-get install -y docker docker-compose net-tools iptables-persistent
    docker network create --ipv6 --subnet "fd00:dead:beef::/48" terminator || true
    touch "APT_STAGE1"
}

mkdir -p /data

docker-compose up -d
