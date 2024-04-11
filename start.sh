#!/bin/sh

set -e

echo "Run migrations"
migrate -path sql/migrations -database "$DB_URL" -verbose up

echo "Start the app"
exec "$@"