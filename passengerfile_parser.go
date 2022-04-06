package passenger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/paketo-buildpacks/packit/v2/fs"
)

type PassengerfileParser struct{}

func NewPassengerfileParser() PassengerfileParser {
	return PassengerfileParser{}
}

type Passengerfile struct {
	Port *int `json:"port,omitempty"`
}

func (p PassengerfileParser) ParsePort(path string, defaultPort int) (int, error) {
	exists, err := fs.Exists(path)
	if err != nil {
		return 0, fmt.Errorf("failed to determine if Passengerfile exists: %w", err)
	}

	if !exists {
		return defaultPort, nil
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("failed to read Passengerfile: %w", err)
	}

	passengerfile := Passengerfile{}
	err = json.Unmarshal(content, &passengerfile)
	if err != nil {
		return 0, fmt.Errorf("failed to parse Passengerfile: %w", err)
	}

	if passengerfile.Port == nil {
		return defaultPort, nil
	}

	return *passengerfile.Port, nil
}
