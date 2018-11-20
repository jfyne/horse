#!/bin/bash

docker run \
    --rm \
    --name horse-postgres \
    --mount source=horse_data,target=/var/lib/postgresql/data \
    -e POSTGRES_PASSWORD=max \
    -e POSTGRES_DB=horse \
    -p 5432:5432 \
    postgres
