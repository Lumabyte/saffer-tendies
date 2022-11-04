#!/usr/bin/env bash

set -e
set -u
set -o pipefail

JSON_DIR="${DATA_DIR}/json"
CSV_DIR="${DATA_DIR}/csv"

main() {
    # Get a list of json files with the full path to them
    mapfile -t json_files < <(find "${JSON_DIR}" -type f -name "*.json")
    # Remove the path, get just the names and .json
#    mapfile -t json_filenames < <(printf "%s\n" "${json_files[@]##*/}")
    # Remove .json to get just the names
#    mapfile -r json_filenames < <(printf "%s\n" "${json_filenames[@]%.json}")

    for json_file in "${json_files[@]}"; do
        json_filename=$(basename "${json_file}" .json)
        printf "%s\n" "Converting ${json_file} into ${CSV_DIR}/${json_filename}.csv"
        jq -r \
            '(.[0] | del(.sd) | keys_unsorted) as $keys | $keys, map([.[ $keys[] ]])[] | @csv' > \
            "${CSV_DIR}/${json_filename}.csv" "${json_file}"

        printf "%s\n" "Converting ${json_file} into ${CSV_DIR}/${json_filename}_sd.csv"
        jq -r \
            '["supportDocumentID","fileName", "extension", "tendersID", "active", "updatedBy", "dateModified", "tenders"], (.[].sd | (.[] | map(.))) | @csv' > \
            "${CSV_DIR}/${json_filename}_sd.csv" "${json_file}"
    done
}

main "${@}"

