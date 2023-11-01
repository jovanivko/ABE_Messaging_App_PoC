package client

import (
	"abeProofOfConcept/internal/logger"
	"abeProofOfConcept/pkg/routes"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"
)

const baseURL = "http://localhost:8080" // Adjust this based on where your server runs

var client *http.Client
var sessionCookie string

func NewHttpClient() {
	jar, _ := cookiejar.New(nil)
	client = &http.Client{Timeout: 20 * time.Second, Jar: jar}
}

func PostRegister(data *routes.RegisterReq) (*http.Response, error) {
	return post("/register", data)
}

func PostLogin(data *routes.LoginReq) (*http.Response, error) {
	return post("/login", data)
}

func PostLogout() (*http.Response, error) {
	return post("/logout", nil)
}

func PostMessage(data *routes.MessageReq) (*http.Response, error) {
	return post("/message", data)
}

func PostFragment(data *routes.FragmentReq) (*http.Response, error) {
	return post("/fragment", data)
}

func GetMailbox() (*http.Response, error) {
	return get("/mailbox")
}

func GetProfile(email string) (*http.Response, error) {
	return get("/profile/" + email)
}

func GetAllEmails() (*http.Response, error) {
	return get("/emails")
}

func post(endpoint string, data interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if sessionCookie != "" {
		req.Header.Set("Cookie", sessionCookie)
	}

	//logger.Logger.Printf("Sending Request to %s with headers: %v\n", baseURL+endpoint, req.Header)
	logger.Logger.Printf("Sending Post Request to %s\n", baseURL+endpoint)
	return client.Do(req)
}

func get(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("GET", baseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if sessionCookie != "" {
		req.Header.Set("Cookie", sessionCookie)
	}

	//logger.Logger.Printf("Sending Request to %s with headers: %v\n", baseURL+endpoint, req.Header)
	logger.Logger.Printf("Sending Get Request to %s\n", baseURL+endpoint)
	return client.Do(req)
}

// Utility functions to read and parse the response from the server
func ReadResponse(response *http.Response, out interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Logger.Errorln("Recovered from panic in ReadResponse:", r)

			err = fmt.Errorf("Recovered from panic in ReadResponse")
		}
	}()
	//logger.Logger.Printf("Received Response with Status: %s and headers: %v\n", response.Status, response.Header)
	logger.Logger.Printf("Received Response with Status: %s", response.Status)
	// Check if the Set-Cookie header is present in the response
	//logger.Logger.Printf("Set-Cookie header: %v\n", response.Header.Get("Set-Cookie"))
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&errorResp)
		if err != nil {
			return err
		}
		return errors.New(errorResp.Message + ": " + errorResp.Error)
	}
	if sc := response.Header.Get("Set-Cookie"); sc != "" {
		sessionCookie = sc
	}

	if out != nil {
		if err = json.NewDecoder(response.Body).Decode(out); err != nil {
			return err
		}
	}
	return nil
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
