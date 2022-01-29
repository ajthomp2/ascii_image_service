# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY main.go ./
COPY asciiimageservice/ ./asciiimageservice/

RUN go build -o /ascii_art_service

EXPOSE 8080

CMD [ "/ascii_art_service" ]
