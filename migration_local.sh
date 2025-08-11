#!/bin/bash
set -e
sleep 2
goose -dir ./migrations postgres "host=pg-local port=5432 dbname=${POSTGRES_DB} user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} sslmode=disable" up -v
