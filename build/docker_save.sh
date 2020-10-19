#!/bin/sh

docker save -o ~/builds/auth_service.tar auth_service
docker save -o ~/builds/resource_service.tar resource_service
docker save -o ~/builds/fetch_service.tar fetch_service
#docker save -o ~/tag_measurements_web_ng.tar tag_measurements_web_ng
docker save -o ~/builds/clean_service.tar clean_service
docker save -o ~/builds/notify_service.tar notify_service
docker save -o ~/builds/tgbot_service.tar tgbot_service
