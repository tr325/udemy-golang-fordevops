package api

import (
	"encoding/json"
	"fmt"
	"io"
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

func (a API) DoGetRequest(requestUrl string) (Response, error) {

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
	var page Page
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
		var words Words
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
		var occurrences Occurrences
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
