#!/bin/sh
set -e

# Map Railway's PORT to app's SERVER_PORT if not explicitly set
if [ -n "${PORT:-}" ] && [ -z "${SERVER_PORT:-}" ]; then
  export SERVER_PORT="$PORT"
fi

# Ensure migrations path is available inside container
if [ -z "${MIGRATION_PATH:-}" ]; then
  export MIGRATION_PATH="/app/database/migrations"
fi

# Run migrations unless explicitly disabled
if [ "${RUN_MIGRATIONS:-true}" != "false" ]; then
  echo "Running database migrations..."
  /app/migration up
  echo "Migrations complete."
else
  echo "Skipping migrations (RUN_MIGRATIONS=false)."
fi

# Start application
exec /app/main
