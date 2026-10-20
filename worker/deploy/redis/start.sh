#!/bin/sh
set -eu
: "${REDIS_PASSWORD:?Set a strong Redis password}"
exec /usr/local/bin/docker-entrypoint.sh redis-server \
  --bind 0.0.0.0 :: \
  --protected-mode yes \
  --requirepass "$REDIS_PASSWORD" \
  --dir /data \
  --appendonly yes \
  --appendfsync always \
  --save "" \
  --maxmemory "${REDIS_MAXMEMORY:-64mb}" \
  --maxmemory-policy noeviction \
  --auto-aof-rewrite-percentage 100 \
  --auto-aof-rewrite-min-size 16mb
