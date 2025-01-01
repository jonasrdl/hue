package hue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Brigde represents a Hue bridge.
type Bridge struct {
	ID        string
	IPAddress string
	Username  string
	client    *http.Client
}

// NewBridge initializes a new Bridge client.
func NewBridge(ip, username string) *Bridge {
	return &Bridge{
		IPAddress: ip,
		Username:  username,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// request performs an HTTP request to the bridge.
func (b *Bridge) request(method, endpoint string, body interface{}) ([]byte, error) {
	url := fmt.Sprintf("http://%s/api/%s/%s", b.IPAddress, b.Username, endpoint)
	var req *http.Request
	var err error

	if body != nil {
		var bodyReader *bytes.Reader
		switch v := body.(type) {
		case []byte:
			bodyReader = bytes.NewReader(v)
		default:
			data, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("failed to serialize body: %w", err)
			}
			bodyReader = bytes.NewReader(data)
		}

		req, err = http.NewRequest(method, url, bodyReader)
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTTP response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return responseData, fmt.Errorf("unexpected HTTP status code: %d", resp.StatusCode)
	}

	return responseData, nil
}
