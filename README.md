# tenders-csv-db
This is a collection of tender data that has been converted into csv and imported into a SQLite DB for further analysis.

# Prereqisites
- curl
- sqlite3
- jq

# Warning
There are some failures around unique constraints, I haven't looked into it yet!

# Setup
Just run the `bootstrap.sh` script like so :)
```
./bootstrap.sh
```

# HTTP Endpoints

- Advertised tenders:
```
https://www.etenders.gov.za/Home/TenderOpportunities/?status=1
```
- Awarded tenders:
```
https://www.etenders.gov.za/Home/TenderOpportunities/?status=2
```

- Closed tenders:
```
https://www.etenders.gov.za/Home/TenderOpportunities/?status=3
```

- Cancelled tenders:
```
https://www.etenders.gov.za/Home/TenderOpportunities/?status=4
```

# SQLite

- Add headers to output
```
.headers on
```

- Format results as columns
```
.mode column
```
