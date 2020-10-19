#!/bin/bash

export USERNAME="an"
export REMOTE_IP="10.10.10.23"
export PASSWORD=""
echo "Please, setup a SSH password for sending docker images: "
read -r -s PASSWORD
if [ -z "$PASSWORD" ];
then
  echo "PASSWORD can't be empty"
else
  sshpass -p "$PASSWORD" scp -rv ~/builds/auth_service.tar $USERNAME@$REMOTE_IP:/mnt/datasource
  sshpass -p "$PASSWORD" scp -rv ~/builds/fetch_service.tar $USERNAME@$REMOTE_IP:/mnt/datasource
  sshpass -p "$PASSWORD" scp -rv ~/builds/resource_service.tar $USERNAME@$REMOTE_IP:/mnt/datasource
  sshpass -p "$PASSWORD" scp -rv ~/builds/notify_service.tar $USERNAME@$REMOTE_IP:/mnt/datasource
  sshpass -p "$PASSWORD" scp -rv ~/builds/clean_service.tar $USERNAME@$REMOTE_IP:/mnt/datasource
#  sshpass -p "$PASSWORD" scp -rv ~/builds/tag_measurements_web_ng.tar $USERNAME@$REMOTE_IP:/mnt/datasource
  sshpass -p "$PASSWORD" scp -rv ~/builds/tgbot_service.tar $USERNAME@$REMOTE_IP:/mnt/datasource
fi

