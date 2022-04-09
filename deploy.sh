#!/bin/bash
APPLICATION_NAME=go-jwt

CONTAINER_ID=$(docker ps | grep go- | awk '{print $1}')

if [ "$CONTAINER_ID" ]
then
  docker stop $CONTAINER_ID
  docker rm $CONTAINER_ID
fi

# cd /home/ec2-user/go-jwt
# docker-compose build
# docker-compose up -d
