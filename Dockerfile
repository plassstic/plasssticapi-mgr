FROM golang:alpine
LABEL authors="plastictactic"
WORKDIR /deploy
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o mgr_ex 'plassstic.tech/gopkg/plassstic-mgr/cmd'