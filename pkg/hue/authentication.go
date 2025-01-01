package hue

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// AuthResponse represents the response from the bridge when creating a user.
type AuthResponse struct {
	Success map[string]interface{} `json:"success"`
	Error   map[string]interface{} `json:"error"`
}

// AuthenticateWithBridge authenticates with the Hue bridge and returns an API key.
// The user must press the link button on the bridge before calling this function.
func AuthenticateWithBridge(ip string, appName string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Payload to create a new user
	payload := map[string]string{
		"devicetype": appName,
	}
	data, _ := json.Marshal(payload)

	url := fmt.Sprintf("http://%s/api", ip)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var response []AuthResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return "", err
	}

	if len(response) == 0 {
		return "", errors.New("unexpected response from bridge: no data returned")
	}

	if response[0].Error != nil {
		errorDetails, _ := json.Marshal(response[0].Error)
		return "", fmt.Errorf("authentication failed: %s", string(errorDetails))
	}

	if username, ok := response[0].Success["username"].(string); ok {
		return username, nil
	}

	return "", errors.New("unexpected response from bridge: username not found")
}
