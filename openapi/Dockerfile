FROM golang:1.24-bullseye AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY main.go main.go

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o app

FROM ubuntu:22.04

RUN mkdir /app
WORKDIR /app
COPY . .
COPY --from=builder /app/app .
ENTRYPOINT ["./app"]