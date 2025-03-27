package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

// const DAREDB_BASE_URL = "https://127.0.0.1:2605"
const DAREDB_BASE_URL = "http://127.0.0.1:2605"

// Example JWT token (replace with a real token)
const JWT_TOKEN = "YOUR-JWT-TOKEN-HERE"

func addKeyWithPost() {
	log.Println("Making POST request to a database")
	url := fmt.Sprintf("%s/set", DAREDB_BASE_URL)
	log.Printf("URL: %s\n", url)

	reqSet, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(`{"myKey":"myValue"}`)))
	if err != nil {
		log.Println("Error creating set request:", err)
		return
	}

	reqSet.Header.Set("Content-Type", "application/json")
	reqSet.Header.Set("Authorization", JWT_TOKEN)
	resp, err := http.DefaultClient.Do(reqSet)
	defer resp.Body.Close()
	if err != nil {
		log.Println("Error while inserting item:", err)
		return
	}

	body := new(bytes.Buffer)
	body.ReadFrom(resp.Body)
	log.Printf("HTTP Code: %d\n", resp.StatusCode)
	log.Printf("content: %s\n\n", body.String())
}

func retrieveKeyWithGet() {
	log.Println("Making GET request to a database")
	keyToRetrieve := "myKey"
	url := fmt.Sprintf("%s/get/%s", DAREDB_BASE_URL, keyToRetrieve)
	log.Printf("URL: %s\n", url)

	reqGet, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error creating get request:", err)
		return
	}
	reqGet.Header.Set("Authorization", JWT_TOKEN)
	resp, err := http.DefaultClient.Do(reqGet)
	if err != nil {
		log.Println("Error while retrieving item:", err)
		return
	}
	defer resp.Body.Close()

	body := new(bytes.Buffer)
	body.ReadFrom(resp.Body)
	log.Printf("HTTP Code: %d\n", resp.StatusCode)
	log.Printf("content: %s\n\n", body.String())
}

func deleteKeyWithDelete() {
	log.Println("Making DELETE request to a database")
	keyToRetrieve := "myKey"
	url := fmt.Sprintf("%s/delete/%s", DAREDB_BASE_URL, keyToRetrieve)
	log.Printf("URL: %s\n", url)

	reqGet, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Println("Error creating get request:", err)
		return
	}
	reqGet.Header.Set("Authorization", JWT_TOKEN)
	resp, err := http.DefaultClient.Do(reqGet)
	if err != nil {
		log.Println("Error while retrieving item:", err)
		return
	}
	defer resp.Body.Close()

	body := new(bytes.Buffer)
	body.ReadFrom(resp.Body)
	log.Printf("HTTP Code: %d\n", resp.StatusCode)
	log.Printf("content: %s\n\n", body.String())
}

func main() {
	addKeyWithPost()
	retrieveKeyWithGet()
	deleteKeyWithDelete()
}
