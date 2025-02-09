package hue

import (
	"encoding/json"
	"net/http"
)

// Scene represents a scene
type Scene struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Owner           string   `json:"owner"`
	Lights          []string `json:"lights"`
	Type            string   `json:"type"`
	Recycle         bool     `json:"recycle"`
	Locked          bool     `json:"locked"`
	CreationTime    string   `json:"created"`
	LastUpdatedTime string   `json:"lastupdated"`
	bridge          *Bridge
}

// GetScenes fetches all scenes from the bridge
func (b *Bridge) GetScenes() ([]Scene, error) {
	data, err := b.request(http.MethodGet, "scenes", nil)
	if err != nil {
		return nil, err
	}

	var scenes map[string]Scene
	if err := json.Unmarshal(data, &scenes); err != nil {
		return nil, err
	}

	var result []Scene
	for id, scene := range scenes {
		scene.ID = id
		scene.bridge = b
		result = append(result, scene)
	}

	return result, nil
}
