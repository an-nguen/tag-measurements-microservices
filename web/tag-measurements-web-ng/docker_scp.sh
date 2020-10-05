#!/bin/bash

export USERNAME=""
export REMOTE_IP=""
export PASSWORD=""
echo "Please, setup a SSH password for sending docker images: "
read -r -s PASSWORD
if [ -z "$PASSWORD" ];
then
  echo "PASSWORD can't be empty"
else
  sshpass -p "$PASSWORD" scp -rv ~/web_thermo_ng_service.tar $USERNAME@$REMOTE_IP:/mnt/datasource
fi

