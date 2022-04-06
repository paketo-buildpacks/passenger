package fakes

import "sync"

type PassengerfileConfigParser struct {
	ParsePortCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Path        string
			DefaultPort int
		}
		Returns struct {
			Int   int
			Error error
		}
		Stub func(string, int) (int, error)
	}
}

func (f *PassengerfileConfigParser) ParsePort(param1 string, param2 int) (int, error) {
	f.ParsePortCall.mutex.Lock()
	defer f.ParsePortCall.mutex.Unlock()
	f.ParsePortCall.CallCount++
	f.ParsePortCall.Receives.Path = param1
	f.ParsePortCall.Receives.DefaultPort = param2
	if f.ParsePortCall.Stub != nil {
		return f.ParsePortCall.Stub(param1, param2)
	}
	return f.ParsePortCall.Returns.Int, f.ParsePortCall.Returns.Error
}
