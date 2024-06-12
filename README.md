# noDare-DB

**noDare-DB** is a project that provides an in-memory database utilizing Redis-inspired hashtables implemented in Go [here](https://github.com/dmarro89/go-redis-hashtable). It offers a lightweight and efficient solution for storing data in memory and accessing it through simple HTTP operations.

## Project Purpose

The primary goal of this project is to offer an in-memory database that leverages hashtables for efficient data storage and retrieval. The Go implementation allows using this database as a component in other Go services or integrating it into applications that require rapid access to in-memory data.

## Running the Database

### Using Docker

To run the database as a Docker image, ensure you have Docker installed on your system. First, navigate to the root directory of your project and execute the following command to build the Docker image:

```bash
docker build -t nodare-db .
```
Once the image is built, you can run the database as a Docker container with the following command:

```bash
docker run -d -p 2420:2420 -p 2240:2240 nodare-db
```

This command will start the database as a Docker container in detached mode exposing ...
... TCP port 2420 of the container to port ```2420``` on your ```localhost```
... UDP port 2240 of the container to port ```2240``` on your ```localhost```

### Using TLS Version in Docker

Build special Docker image, which will generate certificates

```bash
docker build --target nodare-db-tls -f Dockerfile.tls.yml .
```

Once the image is built, you can run the database as a Docker container with the following command:

```bash
docker run -d -e NDB_HOST="0.0.0.0" -p "127.0.0.1:2420:2420" -p "127.0.0.1:2240:2240" -e NDB_PORT=2420 -e NDB_UDP_PORT=2420 -e NDB_SUB_DICKS=1000 -e NDB_TLS_ENABLED="True" -e NDB_TLS_KEY="/app/settings/privkey.pem" -e NDB_TLS_CRT="/app/settings/fullchain.pem" nodare-db-tls
```

Access GET/SET/DEL Commands via API over HTTPS on https://localhost:2420


## How to Use

The in-memory database provides three simple HTTP endpoints to interact with stored data:

### GET /get/{key}

This endpoint retrieves an item from the hashtable using a specific key.

Example usage with cURL:

```bash
curl -X GET http://localhost:2420/get/myKey
```

### SET /set

This endpoint inserts a new item into the hashtable. The request body should contain the key and value of the new item.

Example usage with cURL:

```bash
curl -X POST -d '{"myKey":"myValue"}' http://localhost:2420/set
```

### DELETE /delete/{key}

This endpoint deletes an item from the hashtable using a specific key.

Example usage with cURL:

```bash
curl -X DELETE http://localhost:2420/delete/myKey
```


## Example Usage

Below is a simple example of how to use this database in a Go application:

```go
package main

import (
    "fmt"
    "net/http"
    "bytes"
)

func main() {
    // Example of inserting a new item
    _, err := http.Post("http://localhost:2420/set", "application/json", bytes.NewBuffer([]byte(`{"myKey":"myValue"}`)))
    if err != nil {
        fmt.Println("Error while inserting item:", err)
        return
    }

    // Example of retrieving an item
    resp, err := http.Get("http://localhost:2420/get/myKey")
    if err != nil {
        fmt.Println("Error while retrieving item:", err)
        return
    }
    defer resp.Body.Close()

    // Example of deleting an item
    req, err := http.Get("http://localhost:2420/del/myKey", nil)
    if err != nil {
        fmt.Println("Error while deleting item:", err)
        return
    }
    _, err = http.DefaultClient.Do(req)
    if err != nil {
        fmt.Println("Error while deleting item:", err)
        return
    }
}


