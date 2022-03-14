#!/bin/bash
# Initialization  for the postgres container in the docker compose setup that creates the db and imports data
/setup-db.sh "postgres" "5432" "$POSTGRES_DB" "$DB_USER" "$DB_PASSWORD" "/pokemon-data" "--ignore-host"
