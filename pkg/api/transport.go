package api

import (
	"net/http"
)

type JWTTransport struct {
	transport http.RoundTripper
	token     string
	password  string
	loginURL  string
}

func (t *JWTTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if t.token == "" {
		if t.password != "" {
			token, err := doLoginRequest(http.Client{}, t.loginURL, t.password)
			if err != nil {
				return nil, err
			}
			t.token = token
		}
	}
	if t.token != "" {
		request.Header.Add("Authorization", "Bearer "+t.token)
	}
	return t.transport.RoundTrip(request)
}