package api

import "net/http"

type Options struct {
	Password string
	LoginURL string
}

type apiInterface interface {
	DoGetRequest(requestUrl string) (response, error)
}

type api struct {
	Options Options
	Client  http.Client
}

func New(options Options) apiInterface {
	return api{
		Options: options,
		Client: http.Client{
			Transport: &jwtTransport{
				transport: http.DefaultTransport,
				password:  options.Password,
				loginURL:  options.LoginURL,
			},
		},
	}
}
