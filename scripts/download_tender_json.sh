#!/usr/bin/env bash

set -e
set -u
set -o pipefail

JSON_DIR="${DATA_DIR}/json"

main() {
    printf "%s\n" "Downloading Advertised tenders..."
    curl https://www.etenders.gov.za/Home/TenderOpportunities/?status=1 > "${JSON_DIR}"/etenders_advertised.json
    printf "%s\n" "Downloading Awarded tenders..."
    curl https://www.etenders.gov.za/Home/TenderOpportunities/?status=2 > "${JSON_DIR}"/etenders_awarded.json
    printf "%s\n" "Downloading Closed tenders..."
    curl https://www.etenders.gov.za/Home/TenderOpportunities/?status=3 > "${JSON_DIR}"/etenders_closed.json
    printf "%s\n" "Downloading Cancelled tenders..."
    curl https://www.etenders.gov.za/Home/TenderOpportunities/?status=4 > "${JSON_DIR}"/etenders_cancelled.json
}

main "${@}"
