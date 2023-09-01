package main

import "net/http"

type JWTTransport struct {
	transport http.RoundTripper
	token     string
}

func (t JWTTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if t.token != "" {
		request.Header.Add("Authorization", "Bearer "+t.token)
	}
	return t.transport.RoundTrip(request)
}
