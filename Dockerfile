FROM golang:1.21-bullseye as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o app

FROM ubuntu:22.04
RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/app .
ENTRYPOINT ["./app"]