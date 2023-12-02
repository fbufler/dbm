#!/bin/bash

# start 2 postgres containers
docker run --name postgres1 -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 postgres
docker run --name postgres2 -e POSTGRES_PASSWORD=postgres -d -p 5433:5432 postgres

# start local
go run cmd/main.go local
