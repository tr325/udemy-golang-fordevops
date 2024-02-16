package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

type mockRoundTripper struct {
	output *http.Response
}

func (m mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Note: test we always pass "abc" as our auth token, returning an error
	//  if auth token is anything else (or empty)
	if req.Header.Get("Authorization") != "Bearer abc" {
		return nil, fmt.Errorf("Wrong authorization header: %s", req.Header.Get("Authorization"))
	}
	return m.output, nil
}

func TestRoundTrip(t *testing.T) {
	loginResponse := loginResponse{
		Token: "abc",
	}
	loginResponseBytes, err := json.Marshal(loginResponse)
	if err != nil {
		t.Errorf("marshal error: %s", err)
	}

	myJWTTransport := jwtTransport{
		transport: mockRoundTripper{output: &http.Response{
			StatusCode: 200,
		}},
		httpClient: mockClient{
			postResponse: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(loginResponseBytes)),
			},
		},
		password: "xyz",
	}
	req := &http.Request{
		Header: make(http.Header),
	}

	// Test RoundTrip
	res, err := myJWTTransport.RoundTrip(req)
	if err != nil {
		t.Fatalf("Roundtrip error: %s", err)
	}
	if res.StatusCode != 200 {
		t.Errorf("StatusCode is not 200, got: %d", res.StatusCode)
	}
}
