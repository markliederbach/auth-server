package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	baseURL string = "http://localhost:8080/v1"
)

var (
	loginURL string = fmt.Sprintf("%s/login", baseURL)
)

type LoginRequest struct {
	Username string
	Password string
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func main() {
	// Login
	loginTokens, err := login("erik", "foobar")
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Access Token: %s", loginTokens.AccessToken)
	log.Printf("Refresh Token: %s", loginTokens.RefreshToken)
}

func login(username, password string) (LoginResponse, error) {
	var err error

	client := http.Client{Timeout: time.Second * 3}

	requestBody, err := json.Marshal(LoginRequest{Username: username, Password: password})
	if err != nil {
		return LoginResponse{}, fmt.Errorf("Failed to load JSON request: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, loginURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return LoginResponse{}, fmt.Errorf("Failed to build HTTP Request: %w", err)
	}

	response, err := client.Do(request)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("Failed while calling request: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return LoginResponse{}, fmt.Errorf("Response from %s: %s", loginURL, response.Status)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("Failed to read response body: %w", err)
	}

	var loginResponse LoginResponse
	err = json.Unmarshal(body, &loginResponse)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("Failed to load response JSON: %w", err)
	}

	return loginResponse, nil
}
