## About 

Example in Go to work with DareDB

## Setup

Create `.env` file with the following content (modify credentials)
```
BASE_URL=http://127.0.0.1:2605
#BASE_URL=https://127.0.0.1:2605
DB_USERNAME=<YOUR-USER-HERE>
DB_PASSWORD=<YOUR-PASSWORD-HERE>

```


## Example

* Making simple request to the database using JWT: `daredb_basic_jwt.go`
	```
	go run daredb_basic_jwt.go
	```