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
  sshpass -p "$PASSWORD" scp -rv ~/auth_service.tar $USERNAME@$REMOTE_IP:/mnt/datasource
  sshpass -p "$PASSWORD" scp -rv ~/fetch_service.tar $USERNAME@$REMOTE_IP:/mnt/datasource
  sshpass -p "$PASSWORD" scp -rv ~/resource_service.tar $USERNAME@$REMOTE_IP:/mnt/datasource
  sshpass -p "$PASSWORD" scp -rv ~/notify_service.tar $USERNAME@$REMOTE_IP:/mnt/datasource
  sshpass -p "$PASSWORD" scp -rv ~/clean_service.tar $USERNAME@$REMOTE_IP:/mnt/datasource
  sshpass -p "$PASSWORD" scp -rv ~/web_thermo_ng_service.tar $USERNAME@$REMOTE_IP:/mnt/datasource
  sshpass -p "$PASSWORD" scp -rv ~/tgbot_service.tar $USERNAME@$REMOTE_IP:/mnt/datasource
fi

