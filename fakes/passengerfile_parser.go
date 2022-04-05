package fakes

import (
	"sync"

	"github.com/paketo-buildpacks/passenger"
)

type PassengerfileConfigParser struct {
	ParseCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Path string
		}
		Returns struct {
			Passengerfile passenger.Passengerfile
			Error         error
		}
		Stub func(string) (passenger.Passengerfile, error)
	}
}

func (f *PassengerfileConfigParser) Parse(param1 string) (passenger.Passengerfile, error) {
	f.ParseCall.mutex.Lock()
	defer f.ParseCall.mutex.Unlock()
	f.ParseCall.CallCount++
	f.ParseCall.Receives.Path = param1
	if f.ParseCall.Stub != nil {
		return f.ParseCall.Stub(param1)
	}
	return f.ParseCall.Returns.Passengerfile, f.ParseCall.Returns.Error
}
