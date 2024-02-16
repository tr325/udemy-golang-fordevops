package api

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type response interface {
	GetResponse() string
}

type wordsPage struct {
	page
	words
}

// Using capitals for the type and fields ensures they are exported, and are available to other packages.
// lower case is not exported.
type page struct {
	Name string `json:"page"`
}

type words struct {
	Input string   `json:"input"`
	Words []string `json:"words"`
}

func (w words) GetResponse() string {
	return strings.Join(w.Words, ", ")
}

type occurrences struct {
	Words map[string]int `json:"words"`
}

func (o occurrences) GetResponse() string {
	var response []string
	for word, count := range o.Words {
		response = append(response, fmt.Sprintf("%s (%d)", word, count))
	}
	return strings.Join(response, ", ")
}

func (a api) DoGetRequest(requestUrl string) (response, error) {

	response, err := a.Client.Get(requestUrl)
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
	var page page
	err = json.Unmarshal(body, &page)
	if err != nil {
		return nil, RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("Page unmarshall error: %s", err),
		}
	}

	switch page.Name {
	case "words":
		var words words
		err = json.Unmarshal(body, &words)
		if err != nil {
			return nil, RequestError{
				HTTPCode: response.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("Words unmarshall error: %s", err),
			}
		}
		return words, nil
	case "occurrence":
		var occurrences occurrences
		err = json.Unmarshal(body, &occurrences)
		if err != nil {
			return nil, RequestError{
				HTTPCode: response.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("Occurrences unmarshall error: %s", err),
			}
		}
		return occurrences, nil
	}

	return nil, nil
}
