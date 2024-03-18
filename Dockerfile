FROM golang:1.22.1
FROM ubuntu:latest

ENTRYPOINT ["top", "-b"]

WORKDIR /usr/src/app

COPY . .

RUN go mod tidy