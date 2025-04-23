package main

import (
	"log"
	"io"
	"net/http"
)


func main() {
	
	url := "http://localhost:8080/recipe/random"
	log.Printf("Making request to %v\n", url)
	response, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body) 
	
	log.Printf("Response body: %s", body)
	log.Printf("Status code: %d", response.Status)
}
