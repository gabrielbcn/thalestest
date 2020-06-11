# Super simple dockerfile, feel free to modify

FROM golang:latest

WORKDIR /app

COPY ./ /app

RUN go run thalestest.go
