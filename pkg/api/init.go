package api

import "net/http"

type Options struct {
	Password string
	LoginURL string
}

type APIInterface interface {
	DoGetRequest(requestUrl string) (Response, error)
}

type API struct {
	Options Options
	Client  http.Client
}

func New(options Options) APIInterface {
	return API{
		Options: options,
		Client: http.Client{
			Transport: &JWTTransport{
				transport: http.DefaultTransport,
				password:  options.Password,
				loginURL:  options.LoginURL,
			},
		},
	}
}
