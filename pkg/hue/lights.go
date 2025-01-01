package hue

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Light represents a Philips Hue light resource.
type Light struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	State            State  `json:"state"`
	Type             string `json:"type"`
	ModelID          string `json:"modelid"`
	UniqueID         string `json:"uniqueid"`
	ManufacturerName string `json:"manufacturername"`
	bridge           *Bridge
}

// State represents the state of a light.
type State struct {
	On             *bool           `json:"on,omitempty"`
	Brightness     int             `json:"bri,omitempty"`
	Hue            int             `json:"hue,omitempty"`
	Saturation     int             `json:"sat,omitempty"`
	Effect         string          `json:"effect,omitempty"`
	TransitionTime int             `json:"transitiontime,omitempty"`
	Alert          string          `json:"alert,omitempty"`
	Raw            json.RawMessage `json:"-"` // Allow custom payload overrides
}

// GetLights fetches all lights from the bridge.
func (b *Bridge) GetLights() ([]Light, error) {
	data, err := b.request("GET", "lights", nil)
	if err != nil {
		return nil, err
	}

	var lights map[string]Light
	if err := json.Unmarshal(data, &lights); err != nil {
		return nil, err
	}

	var result []Light
	for id, light := range lights {
		light.ID = id
		result = append(result, light)
	}

	return result, nil
}

// GetLightByID fetches a specific light by its ID from the bridge.
func (b *Bridge) GetLightByID(id string) (*Light, error) {
	data, err := b.request(http.MethodGet, fmt.Sprintf("lights/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var light Light
	if err := json.Unmarshal(data, &light); err != nil {
		return nil, err
	}

	light.ID = id
	light.bridge = b
	return &light, nil
}

// TurnOn turns on the light.
func (l *Light) TurnOn() error {
	on := true
	return l.bridge.SetLightState(l.ID, State{On: &on})
}

// TurnOff turns off the light.
func (l *Light) TurnOff() error {
	on := false
	return l.bridge.SetLightState(l.ID, State{On: &on})
}

// SetBrightness sets the brightness of the light.
func (l *Light) SetBrightness(brightness int) error {
	if brightness < 1 || brightness > 254 {
		return fmt.Errorf("brightness must be between 1 and 254")
	}
	return l.bridge.SetLightState(l.ID, State{Brightness: brightness})
}

// SetColor sets the color of the light by specifying the hue and saturation.
func (l *Light) SetColor(hue, saturation int) error {
	if hue < 0 || hue > 65535 {
		return fmt.Errorf("hue must be between 0 and 65535")
	}
	if saturation < 0 || saturation > 254 {
		return fmt.Errorf("saturation must be between 0 and 254")
	}
	return l.bridge.SetLightState(l.ID, State{Hue: hue, Saturation: saturation})
}

// ToggleLight toggles the on/off state of a light.
func (b *Bridge) ToggleLight(id string, on bool) error {
	payload := map[string]bool{"on": on}
	_, err := b.request("PUT", fmt.Sprintf("lights/%s/state", id), payload)
	return err
}

// SetLightState updates the state of a specified light on the Philips Hue bridge.
// It takes a State struct as input, which specifies the desired properties such as
// power (on/off), brightness, hue, and saturation.
func (b *Bridge) SetLightState(id string, state State) error {
	var data []byte
	var err error

	if state.Raw != nil {
		// Use RawMessage if provided
		data = state.Raw
	} else {
		// Marshal the State struct by default
		data, err = json.Marshal(state)
		if err != nil {
			return fmt.Errorf("failed to serialize payload: %v", err)
		}
	}

	fmt.Printf("Sending payload to light %s: %s\n", id, string(data))

	// Send the HTTP request
	response, err := b.request(http.MethodPut, fmt.Sprintf("lights/%s/state", id), data)
	if err != nil {
		return fmt.Errorf("failed to update light state: %v", err)
	}

	fmt.Printf("Response from bridge: %s\n", string(response))
	return nil
}
