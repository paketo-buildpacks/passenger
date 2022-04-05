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

func (p PassengerfileParser) Parse(path string) (Passengerfile, error) {
	exists, err := fs.Exists(path)
	if err != nil {
		return Passengerfile{}, fmt.Errorf("failed to determine if Passengerfile exists: %w", err)
	}

	if !exists {
		return Passengerfile{}, nil
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return Passengerfile{}, fmt.Errorf("failed to read Passengerfile: %w", err)
	}

	passengerfile := Passengerfile{}
	err = json.Unmarshal(content, &passengerfile)
	if err != nil {
		return Passengerfile{}, fmt.Errorf("failed to parse Passengerfile: %w", err)
	}

	return passengerfile, nil
}
