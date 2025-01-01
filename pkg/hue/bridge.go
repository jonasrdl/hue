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

// BridgeConfig represents the configuration of the Philips Hue bridge.
type BridgeConfig struct {
	Name             string `json:"name"`
	ZigBeeChannel    int    `json:"zigbeechannel"`
	BridgeID         string `json:"bridgeid"`
	MacAddress       string `json:"mac"`
	DHCPActive       bool   `json:"dhcp"`
	IPAddress        string `json:"ipaddress"`
	Netmask          string `json:"netmask"`
	Gateway          string `json:"gateway"`
	ProxyAddress     string `json:"proxyaddress"`
	ProxyPort        int    `json:"proxyport"`
	UTC              string `json:"utc"`
	Localtime        string `json:"localtime"`
	Timezone         string `json:"timezone"`
	ModelID          string `json:"modelid"`
	DatastoreVersion string `json:"datastoreversion"`
	SoftwareVersion  string `json:"swversion"`
	APIVersion       string `json:"apiversion"`
	SoftwareUpdate2  struct {
		CheckForUpdate bool   `json:"checkforupdate"`
		LastChange     string `json:"lastchange"`
		Bridge         struct {
			State       string `json:"state"`
			LastInstall string `json:"lastinstall"`
		} `json:"bridge"`
		State       string `json:"state"`
		AutoInstall struct {
			UpdateTime string `json:"updatetime"`
			On         bool   `json:"on"`
		}
	} `json:"swupdate2"`
	LinkButton       bool   `json:"linkbutton"`
	PortalServices   bool   `json:"portalservices"`
	AnalyticsConsent bool   `json:"analyticsconsent"`
	PortalConnection string `json:"portalconnection"`
	PortalState      struct {
		SignedOn      bool   `json:"signedon"`
		Incoming      bool   `json:"incoming"`
		Outgoing      bool   `json:"outgoing"`
		Communication string `json:"communication"`
	} `json:"portalstate"`
	InternetServices struct {
		Internet       string `json:"internet"`
		RemoteAccess   string `json:"remoteaccess"`
		Time           string `json:"time"`
		SoftwareUpdate string `json:"swupdate"`
	} `json:"internetservices"`
	Factorynew       bool   `json:"factorynew"`
	ReplacesBridgeID string `json:"replacesbridgeid"`
	StarterKitID     string `json:"starterkitid"`
	Backup           struct {
		Status    string `json:"status"`
		ErrorCode int    `json:"errorcode"`
	} `json:"backup"`
	Whitelist map[string]struct {
		LastUseDate string `json:"last use date"`
		CreateDate  string `json:"create date"`
		Name        string `json:"name"`
	} `json:"whitelist"`
}

// GetConfig fetches the bridge's configuration, including software version and connected devices.
func (b *Bridge) GetConfig() (*BridgeConfig, error) {
	data, err := b.request("GET", "config", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get bridge config: %w", err)
	}

	fmt.Println(string(data))

	var config BridgeConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bridge config: %w", err)
	}

	return &config, nil
}
