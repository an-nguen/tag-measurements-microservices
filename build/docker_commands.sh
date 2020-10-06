#!/bin/sh

# change to directory with docker files
cd "$(dirname "$0")"/..

# start to build containers
docker build -t fetch_service -f build/FetchServiceProd.Dockerfile .
docker build -t fetch_service -f build/FetchServiceProd.Dockerfile .
docker build -t clean_service -f build/CleanServiceProd.Dockerfile .
docker build -t clean_service -f build/CleanServiceProd.Dockerfile .
docker build -t resource_service -f build/ResourceServiceProd.Dockerfile .
docker build -t resource_service -f build/ResourceServiceProd.Dockerfile .
docker build -t notify_service -f build/NotifyServiceProd.Dockerfile .
docker build -t notify_service -f build/NotifyServiceProd.Dockerfile .
docker build -t auth_service -f build/AuthServiceProd.Dockerfile .
docker build -t auth_service -f build/AuthServiceProd.Dockerfile .
docker build -t tgbot_service -f build/TgBotServiceProd.Dockerfile .
docker build -t tgbot_service -f build/TgBotServiceProd.Dockerfile .
cd "$(dirname "$0")"/../web/tag-measurements-web-ng/
./build_save.sh

echo 'How to run fetch_service: docker run -it --add-host=database:172.17.0.1 --name fetch_service -d fetch_service'
echo 'How to run clean_service: docker run -it --add-host=database:172.17.0.1 --name clean_service -d clean_service'
echo 'How to run notify_service: docker run -it --add-host=database:172.17.0.1 --name notify_service -d notify_service'
echo 'How to run resource_service: docker run -it --add-host=database:172.17.0.1 -p 10100:10100 --name resource_service -d resource_service'
echo 'How to run auth_service: docker run -it --add-host=database:172.17.0.1 -p 10120:10120 --name auth_service -d auth_service'
echo 'How to run tgbot_service: docker run -it  --add-host=database:172.17.0.1  --name tgbot_service -d tgbot_service'
