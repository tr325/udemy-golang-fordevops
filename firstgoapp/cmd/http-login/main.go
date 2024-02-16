package main

import (
	"flag"
	"fmt"
	"helloworld/pkg/api"
	"net/url"
	"os"
)

func main() {
	var (
		requestUrl string
		password   string
		parsedUrl  *url.URL
		err        error
	)

	flag.StringVar(&requestUrl, "url", "", "URL to access")
	flag.StringVar(&password, "password", "", "use a password to access this API")
	flag.Parse()

	if parsedUrl, err = url.ParseRequestURI(requestUrl); err != nil {
		fmt.Printf("URL is invalid: %s\n", err)
		flag.Usage()
		os.Exit(1)
	}

	apiInstance := api.New(api.Options{
		Password: password,
		LoginURL: parsedUrl.Scheme + "://" + parsedUrl.Host + "/login",
	})

	res, err := apiInstance.DoGetRequest(parsedUrl.String())
	if err != nil {
		if requestErr, ok := err.(api.RequestError); ok {
			fmt.Printf("Error: %s (HTTPCode: %d, Body: %s)\n", requestErr.Err, requestErr.HTTPCode, requestErr.Body)
			os.Exit(1)
		}
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if res == nil {
		fmt.Println("No response")
		os.Exit(1)
	}

	fmt.Printf("Response:\n%s\n", res.GetResponse())
}
