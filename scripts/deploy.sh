#!/bin/bash
# Build a go binary to update the lambda function in AWS
timestamp=$(date "+%F_%T")
cd function
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o ../build/bootstrap main.go
cd ../
cd build
zipName=deploy_$timestamp.zip
zip $zipName bootstrap 
aws lambda update-function-code --function-name  stori-app --zip-file fileb://$zipName