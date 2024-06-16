# Dare-DB

**GDare-DB** is a project that provides an in-memory database utilizing Redis-inspired hashtables implemented in Go [here](https://github.com/dmarro89/go-redis-hashtable). It offers a lightweight and efficient solution for storing data in memory and accessing it through simple HTTP operations.

## Project Purpose

The primary goal of this project is to offer an in-memory database that leverages hashtables for efficient data storage and retrieval. The Go implementation allows using this database as a component in other Go services or integrating it into applications that require rapid access to in-memory data.

## Running the Database

### Using Docker

To run the database as a Docker image, ensure you have Docker installed on your system. First, navigate to the root directory of your project and execute the following command to build the Docker image:

```bash
docker build -t dare-db .
```
Once the image is built, you can run the database as a Docker container with the following command:

```bash
docker run -d -p 2605:2605 dare-db
```

This command will start the database as a Docker container in detached mode, exposing port 2605 of the container to port ```2605``` on your ```localhost```.

### Using TLS Version in Docker

Build special Docker image, which will generate certificates

```bash
docker build -t dare-db-tls -f Dockerfile.tls.yml .
```

Once the image is built, you can run the database as a Docker container with the following command:

```bash
docker run -d -p 2605:2605 dare-db-tls
```

Access API over HTTPS on https://localhost:2605


## How to Use

The in-memory database provides three simple HTTP endpoints to interact with stored data:

### GET /get/{key}

This endpoint retrieves an item from the hashtable using a specific key.

Example usage with cURL:

```bash
curl -X GET http://localhost:2605/get/myKey
```

### SET /set

This endpoint inserts a new item into the hashtable. The request body should contain the key and value of the new item.

Example usage with cURL:

```bash
curl -X POST -d '{"myKey":"myValue"}' http://localhost:2605/set
```

### DELETE /delete/{key}

This endpoint deletes an item from the hashtable using a specific key.

Example usage with cURL:

```bash
curl -X DELETE http://localhost:2605/delete/myKey
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
    _, err := http.Post("http://localhost:2605/set", "application/json", bytes.NewBuffer([]byte(`{"myKey":"myValue"}`)))
    if err != nil {
        fmt.Println("Error while inserting item:", err)
        return
    }

    // Example of retrieving an item
    resp, err := http.Get("http://localhost:2605/get/myKey")
    if err != nil {
        fmt.Println("Error while retrieving item:", err)
        return
    }
    defer resp.Body.Close()

    // Example of deleting an item
    req, err := http.NewRequest("DELETE", "http://localhost:2605/delete/myKey", nil)
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
```

## How to Contribute

All sorts of contributions to this project!  Here's how you can get involved:

* *Found a bug?* Let us know! Open an [issue](https://github.com/dmarro89/dare-db/issues) and briefly describe the problem.
* *Have a great idea for a new feature?* Open an [issue](https://github.com/dmarro89/dare-db/issues) to discuss it. If you'd like to implementing it yourself, you can assign this issue to yourself and create a pull request once the code/improvement/fix is ready.
* *Want to talk about something related to the project?* [Discussion threads](https://github.com/dmarro89/dare-db/discussions) are the perfect place to brainstorm ideas


Here is how you could add your new ```code/improvement/fix``` with a *pull request*:

1. Fork the repository (e.g., latest changes must be in ```develop``` branch)
    ```
    git clone -b develop https://github.com/dmarro89/dare-db
    ```
2. Create a new branch for your feature. Use number of a newly created issue and keywords (e.g., ```10-implement-feature-ABC```)
    ```
    git checkout -b 10-implement-feature-ABC
    ```
3. Add changes to the branch
    ```
    git add .
    ```
4. Commit your changes 
    ```
    git commit -am 'add new feature ABC'
    ```
5. Push to the branch
    ```
    git push origin 10-implement-feature-ABC
    ```
6. Open a pull request based on a new branch
7. Provide a short notice in the pull request according to the following template:
    + Added: ...
    + Changed: ...
    + Fixed: ...
    + Dependencies: ...