#!/bin/sh

host="$1"
shift
port="$1"
shift
cmd="$@"

# Ожидание до тех пор, пока сервис на host:port не станет доступен
until nc -z -v -w30 $host $port; do
  echo "Waiting for database at $host:$port..."
  sleep 1
done

echo "Database is ready"
exec $cmd
