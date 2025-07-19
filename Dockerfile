FROM ubuntu:24.04

LABEL maintainer="juliotorres"

RUN mkdir -p /app
WORKDIR /app

COPY bin/odoo-executor .

RUN chmod +x /app/odoo-executor

USER ubuntu

EXPOSE 4080

ENTRYPOINT ["/app/odoo-executor"]