package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type loginRequest struct {
	Password string `json:"password"`
}
type loginResponse struct {
	Token string `json:"token"`
}

func doLoginRequest(client http.Client, loginUrl, password string) (string, error) {
	loginRequest := loginRequest{
		Password: password,
	}

	body, err := json.Marshal(loginRequest)
	if err != nil {
		return "", fmt.Errorf("marshall error: %s", err)
	}

	response, err := client.Post(loginUrl, "application/json", bytes.NewBuffer(body))

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
	if !json.Valid(responseBody) {
		return "", RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("No valid JSON returned"),
		}
	}
	var loginResponse loginResponse
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
