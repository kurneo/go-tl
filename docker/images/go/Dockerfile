FROM golang:1.22-alpine

MAINTAINER GiangNguyen <giangnguyen.neko.130@gmail.com>

WORKDIR /var/www

RUN apk update && apk upgrade

RUN apk add --no-cache curl make ffmpeg aws-cli

RUN go install -tags 'postgres,mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN go install github.com/google/wire/cmd/wire@latest

RUN go install github.com/cosmtrek/air@latest

ENTRYPOINT cd /var/www && air