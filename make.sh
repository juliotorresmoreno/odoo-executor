#!/bin/bash

mkdir -p bin
go build -o bin/odoo-executor main.go

docker build -t jliotorresmoreno/odoo-executor:v1.0.0 .

docker push jliotorresmoreno/odoo-executor:v1.0.0
