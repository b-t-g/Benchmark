#!/bin/bash
set -uEe

THIS_SCRIPT_DIR="$(cd $(dirname "${BASH_SOURCE[0]}") && pwd)"
CPU_USAGE_DIR="$(cd "${THIS_SCRIPT_DIR}/../cpu_usage"&& pwd)"

USERNAME="${USERNAME:-postgres}"
PASSWORD="${PASSWORD:-example}"
DB_HOSTNAME="${DB_HOSTNAME:-db}"
PORT="${PORT:-5432}"
DATABASE="${DATABASE:-postgres}"
SSL_MODE="${SSL_MODE:-disable}"

wait_for_db() {
	while ! nc -zv "${DB_HOSTNAME}" 5432; do
		sleep 5
	done
}

# Since we create tables and import data at the same step, use the existence of the table
# as a proxy for whether we also need to import data.
data_exists() {
	if psql "postgres://${USERNAME}:${PASSWORD}@${DB_HOSTNAME}:${PORT}?sslmode=${SSL_MODE}" -lqt | cut -d \| -f 1 | grep -qw "homework"; then
		return 0
	else
		echo "$contains_homework"
		return 1
	fi
}

create_tables() {
	psql "postgres://${USERNAME}:${PASSWORD}@${DB_HOSTNAME}:${PORT}?sslmode=${SSL_MODE}" < "${CPU_USAGE_DIR}/cpu_usage.sql"
}

import_data() {
	psql "postgres://${USERNAME}:${PASSWORD}@${DB_HOSTNAME}:${PORT}/homework?sslmode=${SSL_MODE}" -c "\COPY cpu_usage FROM ${CPU_USAGE_DIR}/cpu_usage.csv CSV HEADER"
}

wait_for_db

if ! data_exists; then
	create_tables
	import_data
else
	echo "Table already exists. Skipping table creation and importing data."
fi
