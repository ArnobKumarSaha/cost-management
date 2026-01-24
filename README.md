
## Cost Management
First create a `.env` file with content
```bash
MONGODB_URI="mongodb+srv://<>:<>@arnob.0b4vj7e.mongodb.net/?retryWrites=true&w=majority"
```

### Import
```bash
go run *.go --mode=import --file=january.txt
```

### Report
```bash
go run *.go --mode=report --period=monthly --month=1

Monthly Cost Sheet for January 2026
-----------------------------
Transport: 80.00
Snacks: 20.00
Bazar: 60.00
House Rent: 23000.00
Other: 7500.00
-----------------------------
Total: 30660.00
```