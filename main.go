package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Using capitals for the type and fields ensures they are exported, and are available to other packages.
// lower case is not exported.
type Page struct {
	Name string `json:"page"`
}

type Words struct {
	Input string   `json:"input"`
	Words []string `json:"words"`
}

type Occurrences struct {
	Words map[string]int `json:"words"`
}

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
	if response.StatusCode != 200 {
		fmt.Printf("Invalid output\nStatus: %d\nBody: %s\n", response.StatusCode, body)
		os.Exit(1)
	}

	// Note: we are parsing the JSON _partially_ here, since we know that both endpoints include the
	// `page` key, and then parsing the rest of the payload separately in the switch statement below
	var page Page
	err = json.Unmarshal(body, &page)
	if err != nil {
		log.Fatal(err)
	}

	switch page.Name {
	case "words":
		var words Words
		err = json.Unmarshal(body, &words)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("JSON Parsed:\nPage: %s\nWords: %s\n", page.Name, strings.Join(words.Words, ", "))
	case "occurrence":
		var occurrences Occurrences
		err = json.Unmarshal(body, &occurrences)
		if err != nil {
			log.Fatal(err)
		}
		for word, count := range occurrences.Words {
			fmt.Printf("%s: %d\n", word, count)
		}
	default:
		fmt.Println("Page not found")
	}

}
