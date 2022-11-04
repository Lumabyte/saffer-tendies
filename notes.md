# Notes
Just random notes as I go along.

## etenders endpoints

### Using curl

Just a note that when using curl you need to escape special bash characters, like & but doing `\&`.

### Download JSON data of whole category
https://www.etenders.gov.za/Home/TenderOpportunities/?status=1 - Advertised
https://www.etenders.gov.za/Home/TenderOpportunities/?status=2 - Awarded
https://www.etenders.gov.za/Home/TenderOpportunities/?status=3 - Closed
https://www.etenders.gov.za/Home/TenderOpportunities/?status=4 - Cancelled

Can just hit these endpoints with curl to download it's data
```
curl https://www.etenders.gov.za/Home/TenderOpportunities/?status=1 > etenders_advertised.json
```

### Query by tender name (keyword)
https://www.etenders.gov.za/Home/SearchByTenderName/?tenderName=<keyword>

```
curl https://www.etenders.gov.za/Home/SearchByTenderName/?tenderName=sftp
```

This returns:
```
[
  {
    "id": 38723,
    "name": "PROCUREMENT OF SFTP SOFTWARE"
  }
]
```

### Query by tender id
This ID is the ID assigned by etenders themselves, this is NOT the tender number.

Using the result from querying by keyword, you can then call this endpoint with the id:
```
curl https://www.etenders.gov.za/Home/HomeSearch/?status=0\&searchTerm=38723
```

The above will return:
```
[
  {
    "id": 38723,
    "tender_No": "RFQ:116/22/ICT",
    "type": "Request for Quotation",
    "delivery": "CNR BEATRIX AND PRETORIOUS STREET - ARCADIA - PRETORIA - 0083",
    "department": "South African Social Security Agency",
    "date_Published": "2022-11-02T00:00:00",
    "cbrief": false,
    "cd": "Monday, 07 November 2022 - 11:00",
    "dp": "Wednesday, 02 November 2022",
    "closing_Date": "2022-11-07T11:00:00",
    "brief": "<not available>",
    "compulsory_briefing_session": null,
    "status": "Published",
    "category": "Services: General",
    "description": "PROCUREMENT OF SFTP SOFTWARE",
    "province": "Gauteng",
    "contactPerson": "DINEO",
    "email": "AcquisitionDineo@sassa.gov.za",
    "telephone": "012-400-2154",
    "fax": "N/A",
    "briefingVenue": null,
    "conditions": "N/A",
    "sd": [
      {
        "supportDocumentID": "98ae2adc-af79-4956-a62e-a9329ff91fd5",
        "fileName": "ETENDER SPEC.pdf",
        "extension": ".pdf",
        "tendersID": 38723,
        "active": true,
        "updatedBy": "DineoLe@sassa.gov.za",
        "dateModified": "2022-11-02T09:29:14.1819343",
        "tenders": null
      }
    ],
    "bf": " NO",
    "bc": " NO"
  }
]
```

## Parsing with jq
### Turns all json into csv excluding 'sd' objects
jq -r '(.[0] | del(.sd) | keys_unsorted) as $keys | $keys, map([.[ $keys[] ]])[] | @csv' > etenders_awarded.csv

### Turns all 'sd' json objects into csv - fails on nulls
jq -r'.[].sd | (map(keys) | add | unique) as $cols | map(. as $row | $cols | map($row[.])) as $rows | $cols, $rows[] | @csv'

### Turns all 'sd' json objects into csv ignoring nulls
Creates csv headings again for every 'sd' object
.[].sd | ([.[] | keys[]] | unique) as $keys | $keys, (.[] | [.[$keys[]]]) | @csv

### Creates no csv headings
.[].sd | (.[] | to_entries | map(.value))|@csv
.[].sd | (.[] | map(.))|@csv

### Headers before the parsing
(["supportDocumentID","fileName", "extension", "tendersID", "active", "updatedBy", "dateModified", "tenders"], ...)

### Getting headers manually
'["supportDocumentID","fileName", "extension", "tendersID", "active", "updatedBy", "dateModified", "tenders"], (.[].sd | (.[] | map(.))) | @csv'
jq -r '["supportDocumentID","fileName", "extension", "tendersID", "active", "updatedBy", "dateModified", "tenders"], (.[].sd | (.[] | map(.))) | @csv' > etenders_awarded_sd.csv

## SQL
### Importing a csv
.mode csv
.import etenders_advertised_sd.csv advertised_sd

### Example queries
```
select awarded.id as id, awarded.tender_No as tender_no, awarded.department, awarded.closing_Date as closing_date, awarded.contactPerson as contact_person, awarded.email from awarded inner join awarded_sd on awarded.id=awarded_sd.tendersID
```

```
select awarded.id, awarded_sd_test.tendersID, awarded.tender_No, awarded.department, awarded.contactPerson, awarded.description, awarded.email, awarded_sd_test.fileName
from awarded
inner join awarded_sd_test on awarded.id=awarded_sd_test.tendersID;
```