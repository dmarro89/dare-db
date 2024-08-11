FROM golang:1.22-bullseye as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o app

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y curl netcat

ENV DARE_TLS_ENABLED="false"
RUN mkdir /app
WORKDIR /app

COPY auth/rbac_model.conf /app/auth/rbac_model.conf
COPY auth/rbac_policy.csv /app/auth/rbac_policy.csv

COPY --from=builder /app/app .

EXPOSE 2605

ENTRYPOINT ["./app"]