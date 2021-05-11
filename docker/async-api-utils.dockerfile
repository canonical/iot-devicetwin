FROM ubuntu:20.04

RUN apt update && DEBIAN_FRONTEND=noninteractive apt install -y nodejs npm golang

RUN npm install -g @asyncapi/generator

RUN DEBIAN_FRONTEND=noninteractive apt install -y git

RUN go get github.com/asyncapi/parser-go/...