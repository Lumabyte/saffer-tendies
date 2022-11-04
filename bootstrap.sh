#!/usr/bin/env bash

set -e
set -u
set -o pipefail

PROJ_DIR="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
DATA_DIR="${PROJ_DIR}/data"
SCRIPTS_DIR="${PROJ_DIR}/scripts"
SQL_DIR="${PROJ_DIR}/sql"

export PROJ_DIR
export DATA_DIR
export SCRIPTS_DIR
export SQL_DIR

main() {
    # 1) Download JSON data from etenders.gov.za
    "${SCRIPTS_DIR}"/download_tender_json.sh
    # 2) Convert JSON data into CSV
    "${SCRIPTS_DIR}"/convert_to_csv.sh
    # 3) Create DB and import CSV data
    "${SCRIPTS_DIR}"/sqlite_import_csv.sh
}

main "${@}"