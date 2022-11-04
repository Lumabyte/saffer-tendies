#!/usr/bin/env bash

set -e
set -u
set -o pipefail

CSV_DIR="${DATA_DIR}/csv"
SQL_DB="tenders.db"

main() {
    # Create the sqlite DB with predefined tables
    sqlite3 "${SQL_DB}" < "${SQL_DIR}"/create_tables.sql

    # Get a list of csv files with the full path to them
    mapfile -t csv_files < <(find "${CSV_DIR}" -type f -name "*.csv")

#    sqlite3 tenders.db ".mode csv"
    for csv_file in "${csv_files[@]}"; do
        csv_filename=$(basename "${csv_file}" .csv)
        tablename="${csv_filename##etenders_}"
        # Perform the actual import
        printf "%s\n" "Importing file ${csv_file} into ${tablename}"
        sqlite3 "${SQL_DB}" ".mode csv" ".import ${csv_file} ${tablename}"
    done

}

main "${@}"

