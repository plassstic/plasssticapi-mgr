FROM golang:latest
LABEL authors="plastictactic"
RUN apk update && apk upgrade && \
    apk add --no-cache bash openssh
WORKDIR /deploy
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD ["go", "run", "plassstic.tech/gopkg/golang-manager/cmd"]