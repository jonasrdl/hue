# Hue Go library

## Setup
If you already know the IP of your bridge and have an existing account, you can use the following:
```go
package main

import (
	"github.com/jonasrdl/hue"
)

func main() {
	bridge := hue.NewBridge("192.168.0.130", "username")

	lights, err := bridge.GetLights()
	if err != nil {
		panic(err)
    }
	
	// ...
}
```

It's also possible to discover the bridge IP and authenticate to the bridge (linking button needs to be pressed beforehand):
```go
package main

import (
	"github.com/jonasrdl/hue"
)

func main() {
	bridgeIP, err := hue.DiscoverBridgeIP()
	if err != nil {
		fmt.Println(err)
	}
	username, err := hue.AuthenticateWithBridge(bridgeIP, "your-app-name") // Press the button before
	if err != nil {
		fmt.Println(err)
	}

	hue.NewBridge(username, bridgeIP)
}
```

## Support Matrix

| **Category**           | **Functionality**        | **Description**                                                                       | **Status** |
|------------------------|--------------------------|---------------------------------------------------------------------------------------|------------|
| **Lights**             | `GetLights`              | Fetches all lights connected to the bridge.                                           | ✔️         |
|                        | `GetLightByID`           | Fetches details of a specific light by its ID.                                        | ✔️         |
|                        | `TurnOn`                 | Turns on a specific light.                                                            | ✔️         |
|                        | `TurnOff`                | Turns off a specific light.                                                           | ✔️         |
|                        | `SetLightState`          | Sets the state of a specific light (e.g., brightness, hue, saturation).               | ✔️         |
|                        | `ToggleLight`            | Toggles the on/off state of a specific light.                                         | ✔️         |
|                        | `SetBrightness`          | Adjusts the brightness of a specific light.                                           | ✔️         |
|                        | `SetColor`               | Sets the color of a specific light (e.g., hue, saturation).                           | ✔️         |
| **Authentication**     | `AuthenticateWithBridge` | Authenticates with the Hue bridge and returns an API key                              | ✔️         |
| **Groups**             | `GetGroups`              | Fetches all groups of lights defined on the bridge.                                   | ✔️          |
|                        | `GetGroupByID`           | Fetches details of a specific group by its ID.                                        | ✔️          |
|                        | `SetGroupState`          | Sets the state of a group of lights (e.g., turn on, adjust brightness).               | ✔️          |
| **Scenes**             | `GetScenes`              | Fetches all scenes available on the bridge.                                           | ❌          |
|                        | `ActivateScene`          | Activates a specific scene by ID.                                                     | ❌          |
| **Sensors**            | `GetSensors`             | Fetches all sensors connected to the bridge.                                          | ❌          |
|                        | `GetSensorByID`          | Fetches details of a specific sensor by its ID.                                       | ❌          |
| **Schedules**          | `GetSchedules`           | Fetches all schedules defined on the bridge.                                          | ❌          |
|                        | `CreateSchedule`         | Creates a new schedule to automate actions (e.g., turn on lights at a specific time). | ❌          |
|                        | `DeleteSchedule`         | Deletes a schedule by its ID.                                                         | ❌          |
| **Bridge Information** | `GetConfig`              | Fetches the bridge's configuration, including software version and connected devices. | ✔️         |
|                        | `GetCapabilities`        | Fetches the bridge's capabilities (e.g., maximum number of lights or groups).         | ✔️          |
