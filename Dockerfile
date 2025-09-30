FROM golang:alpine
LABEL authors="plastictactic"
RUN apk update && apk upgrade && \
    apk add --no-cache bash openssh
WORKDIR /deploy
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o mgr_ex 'plassstic.tech/gopkg/golang-manager/cmd'