#!/bin/bash

set -x

echo "start deploy ${USER}"
GOOS=linux go build -o isuumo -v .
for server in isu01 isu02 isu03; do
  echo "deploy to $server"
  ssh -t $server "sudo systemctl stop isuumo.go.service"
  scp ./isuumo $server:/home/isucon/isuumo/webapp/go/isuumo
  rsync -vau ../mysql/ $server:/home/isucon/isuumo/webapp/mysql/
  ssh -t $server "sudo systemctl start isuumo.go.service"
done

echo "finish deploy ${USER}"
