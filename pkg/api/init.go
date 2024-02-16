package api

import (
	"io"
	"net/http"
)

type Options struct {
	Password string
	LoginURL string
}

type ClientInterface interface {
	Get(url string) (resp *http.Response, err error)
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}
type ApiInterface interface {
	DoGetRequest(requestUrl string) (response, error)
}

type api struct {
	Options Options
	Client  ClientInterface
}

func New(options Options) ApiInterface {
	return api{
		Options: options,
		Client: &http.Client{
			Transport: &jwtTransport{
				transport:  http.DefaultTransport,
				password:   options.Password,
				loginURL:   options.LoginURL,
				httpClient: &http.Client{},
			},
		},
	}
}

func NewWithClient(clientInterface ClientInterface, options Options) ApiInterface {
	return api{
		Options: options,
		Client:  clientInterface,
	}
}
