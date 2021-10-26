#!/bin/bash

# =========================(test)=========================
docker-compose run golang go test api/handlers/auth/* -v
docker rm $(docker ps -a -q -f status=exited)
# =========================(test)=========================

# =========================(run)=========================
docker-compose up -d --build
# =========================(run)=========================
