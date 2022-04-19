#!/bin/sh

set -e

echo "start script..."
db_source="${DB_DRIVER}"://"${DB_USER}":"${DB_PASSWORD}"@"${DB_HOST}":"${DB_PORT}"/"${DB_DATABASE}?sslmode=disable"
echo "run db migration for $db_source"
/app/migrate -path /app/migration -database "$db_source" -verbose up

echo "start the app"
exec "$@"