# syntax=docker/dockerfile:1

## Build
FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./


COPY *.go ./

RUN go build -o ./env-on-api

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /env-on-api /env-on-api

EXPOSE 8088

USER nonroot:nonroot

ENTRYPOINT ["/env-on-api"]