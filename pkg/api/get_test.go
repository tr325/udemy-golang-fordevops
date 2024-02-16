package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type mockClient struct {
	Response *http.Response
}

func (m mockClient) Get(url string) (resp *http.Response, err error) {
	return m.Response, nil
}

func TestDoGetRequest(t *testing.T) {
	wordsPage := wordsPage{
		page: page{"words"},
		words: words{
			Input: "abc",
			Words: []string{"a", "b"},
		},
	}
	pageBytes, err := json.Marshal(wordsPage)
	if err != nil {
		t.Errorf("marshal error: %s", err)
	}

	apiInstance := NewWithClient(mockClient{
		Response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(pageBytes)),
		},
	}, Options{})

	response, err := apiInstance.DoGetRequest("http://localhost/words")
	if err != nil {
		t.Errorf("DoGetRequest error: %s", err)
	}
	if response == nil {
		t.Fatalf("response is empty")
	}
	if response.GetResponse() != "a, b" {
		t.Errorf("Unexpected response: %s", response.GetResponse())
	}
}
