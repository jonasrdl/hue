package hue

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Group represents a group of lights
type Group struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Lights []string `json:"lights"`
	bridge *Bridge
}

func (b *Bridge) GetGroups() ([]Group, error) {
	data, err := b.request(http.MethodGet, "groups", nil)
	if err != nil {
		return nil, err
	}

	var groups map[string]Group
	if err := json.Unmarshal(data, &groups); err != nil {
		return nil, err
	}

	var result []Group
	for id, group := range groups {
		group.ID = id
		group.bridge = b
		result = append(result, group)
	}

	return result, nil
}

// GetGroupByID fetches a specific group by its ID
func (b *Bridge) GetGroupByID(id string) (*Group, error) {
	data, err := b.request(http.MethodGet, fmt.Sprintf("groups/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var group Group
	if err := json.Unmarshal(data, &group); err != nil {
		return nil, err
	}

	group.ID = id
	group.bridge = b
	return &group, nil
}

// SetGroupState updates the state of all lights in a specified group on the Philips Hue bridge.
// It takes a State struct as input, which specifies the desired properties such as
// power (on/off), brightness, hue, and saturation.
func (b *Bridge) SetGroupState(id string, state State) error {
	var data []byte
	var err error

	if state.Raw != nil {
		data = state.Raw
	} else {
		data, err = json.Marshal(state)
		if err != nil {
			return fmt.Errorf("failed to marshal state: %w", err)
		}
	}

	_, err = b.request(http.MethodPut, fmt.Sprintf("groups/%s/action", id), data)
	if err != nil {
		return fmt.Errorf("failed to set state: %w", err)
	}

	return nil
}
