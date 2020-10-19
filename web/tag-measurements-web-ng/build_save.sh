#!/bin/bash

cd "$(dirname "$0")"/.

rm -rf ./dist
ng build --prod
docker build -t tag_measurements_web_ng -f Dockerfile .
docker save tag_measurements_web_ng -o ~/builds/tag_measurements_web_ng.tar


