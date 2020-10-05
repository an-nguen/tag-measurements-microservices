#!/bin/sh

docker save -o ~/auth_service.tar auth_service
docker save -o ~/resource_service.tar resource_service
docker save -o ~/fetch_service.tar fetch_service
docker save -o ~/web_thermo_ng_service.tar web_thermo_ng_service
docker save -o ~/clean_service.tar clean_service
docker save -o ~/notify_service.tar notify_service
docker save -o ~/tgbot_service.tar tgbot_service
