#!/bin/bash

set -e

cd "$(dirname $0)/.."
a=ogiekako@instance-1.asia-northeast1-a.laundry-209013
scp -r laundry $a:~/src/github.com/ogiekako/
ssh $a -t 'sudo systemctl restart laundry'

echo "open http://35.200.14.198:8080"
