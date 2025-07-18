FROM ubuntu:24.04

LABEL maintainer="juliotorres"

RUN apt-get update && \
    apt-get install -y \
    ca-certificates \
    golang-go \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /app
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o /app/main
RUN chown -R ubuntu:ubuntu /app && chmod 755 /app
RUN chmod +x /app/main

USER ubuntu

