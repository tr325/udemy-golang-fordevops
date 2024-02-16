package api

import (
	"net/http"
)

type jwtTransport struct {
	transport  http.RoundTripper
	token      string
	password   string
	loginURL   string
	httpClient ClientInterface
}

func (t *jwtTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if t.token == "" {
		if t.password != "" {
			token, err := doLoginRequest(t.httpClient, t.loginURL, t.password)
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
