#!/bin/sh
#docker run -d -p 27017:27017 --name=mongo-db-container mongo:latest

# Build a container image of docker to run localy the lambda function using lambda rie, see README to know how it works
#docker build --platform linux/amd64 -t stori-app-function:stori-app .
#docker run -v ~/.aws:/root/.aws --env-file ./function/.env -d -p 9000:8080 --entrypoint /usr/local/bin/aws-lambda-rie stori-app-function:stori-app ./bootstrap
docker-compose up --build 