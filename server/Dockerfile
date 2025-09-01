FROM golang:1.25.0-alpine

RUN apk update
RUN apk upgrade

RUN go install github.com/air-verse/air@latest

WORKDIR /app

CMD ["air"]
