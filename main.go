package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage: ./http-get <url>")
		os.Exit(1)
	}

	if _, err := url.ParseRequestURI(args[1]); err != nil {
		fmt.Printf("URL is invalid: %s\n", err)
		os.Exit(1)
	}

	response, err := http.Get(args[1])
	if err != nil {
		log.Fatal(err)
	}
	// defer keyword ensures this is run at the end of this method
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Status: %d\nBody: %s\n", response.StatusCode, body)
}
