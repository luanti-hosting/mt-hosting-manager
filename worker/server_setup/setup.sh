#!/bin/sh
set -e
cd `dirname $0`

docker network create "network-{{.ServerShortID}}" || true
docker compose pull
docker compose up -d --remove-orphans
