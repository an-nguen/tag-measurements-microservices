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
  sshpass -p "$PASSWORD" scp -rv ~/builds/tag_measurements_web_ng.tar $USERNAME@$REMOTE_IP:/mnt/datasource
fi

