#!/bin/bash

docker build -t jliotorresmoreno/odoo-executor:v1.0.0 .

docker push jliotorresmoreno/odoo-executor:v1.0.0
