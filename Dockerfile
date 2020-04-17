FROM golang:buster AS development

#RUN mkdir /app
WORKDIR /go/src/gochat

LABEL maintainer=""

# Here we would install any software dependencies first before copying all source files
RUN apt-get update && apt install git && go get github.com/stretchr/gomniauth github.com/gorilla/websocket 

VOLUME /go/src/gochat

# COPY . .

# https://docs.docker.com/develop/develop-images/multistage-build/
#FROM golang:buster --platform=$OS AS debug

#FROM golang:1-alpine --platform=$OS AS builder

#FROM alpine:latest --platform=$OS AS prod