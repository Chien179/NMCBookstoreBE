#!/bin/sh

set -e

echo "run db migrations"

# echo "vm.overcommit_memory = 1" | tee /etc/sysctl.d/nextcloud-aio-memory-overcommit.conf

echo "start the app"
exec "$@"