package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LoginRequest struct {
	Password string `json:"password"`
}
type LoginResponse struct {
	Token string `json:"token"`
}

func doLoginRequest(loginUrl, password string) (string, error) {
	loginRequest := LoginRequest{
		Password: password,
	}

	body, err := json.Marshal(loginRequest)
	if err != nil {
		return "", fmt.Errorf("marshall error: %s", err)
	}

	response, err := http.Post(loginUrl, "application/json", bytes.NewBuffer(body))

	if err != nil {
		return "", fmt.Errorf("http.Post error: %s", err)
	}
	// defer keyword ensures this is run at the end of this method
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("ReadAll error: %s", err)
	}
	if response.StatusCode != 200 {
		return "", fmt.Errorf("Invalid output\nStatus: %d\nBody: %s", response.StatusCode, responseBody)
	}

	var loginResponse LoginResponse
	err = json.Unmarshal(responseBody, &loginResponse)
	if err != nil {
		return "", RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(responseBody),
			Err:      fmt.Sprintf("LoginResponse unmarshall error: %s", err),
		}
	}

	if loginResponse.Token == "" {
		return "", RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(responseBody),
			Err:      "Empty token replied",
		}
	}

	return loginResponse.Token, nil
}