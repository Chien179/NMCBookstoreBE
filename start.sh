#!/bin/sh

set -e

echo "run db migrations"

echo "start the app"
exec "$@"