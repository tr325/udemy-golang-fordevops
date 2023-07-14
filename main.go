package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Response interface {
	GetResponse() string
}

// Using capitals for the type and fields ensures they are exported, and are available to other packages.
// lower case is not exported.
type Page struct {
	Name string `json:"page"`
}

type Words struct {
	Input string   `json:"input"`
	Words []string `json:"words"`
}

func (w Words) GetResponse() string {
	return strings.Join(w.Words, ", ")
}

type Occurrences struct {
	Words map[string]int `json:"words"`
}

func (o Occurrences) GetResponse() string {
	var response []string
	for word, count := range o.Words {
		response = append(response, fmt.Sprintf("%s (%d)", word, count))
	}
	return strings.Join(response, ", ")
}

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage: ./http-get <url>")
		os.Exit(1)
	}

	res, err := doRequest(args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if res == nil {
		fmt.Println("No response")
		os.Exit(1)
	}

	fmt.Printf("Response:\n%s\n", res.GetResponse())
}

func doRequest(requestUrl string) (Response, error) {
	if _, err := url.ParseRequestURI(requestUrl); err != nil {
		return nil, fmt.Errorf("URL is invalid: %s", err)
	}

	response, err := http.Get(requestUrl)
	if err != nil {
		return nil, fmt.Errorf("http.Get error: %s", err)
	}
	// defer keyword ensures this is run at the end of this method
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("ReadAll error: %s", err)
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Invalid output\nStatus: %d\nBody: %s", response.StatusCode, body)
	}

	// Note: we are parsing the JSON _partially_ here, since we know that both endpoints include the
	// `page` key, and then parsing the rest of the payload separately in the switch statement below
	var page Page
	err = json.Unmarshal(body, &page)
	if err != nil {
		return nil, fmt.Errorf("Unmarshall error: %s", err)
	}

	switch page.Name {
	case "words":
		var words Words
		err = json.Unmarshal(body, &words)
		if err != nil {
			return nil, fmt.Errorf("Unmarshall error: %s", err)
		}
		return words, nil
	case "occurrence":
		var occurrences Occurrences
		err = json.Unmarshal(body, &occurrences)
		if err != nil {
			return nil, fmt.Errorf("Unmarshall error: %s", err)
		}
		return occurrences, nil
	}

	return nil, nil
}
