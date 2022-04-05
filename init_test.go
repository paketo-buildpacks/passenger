package passenger_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitPassenger(t *testing.T) {
	suite := spec.New("passenger", spec.Report(report.Terminal{}), spec.Parallel())
	suite("Build", testBuild)
	suite("Detect", testDetect)
	suite("GemfileParser", testGemfileParser)
	suite("PassengerfileParser", testPassengerfileParser)
	suite.Run(t)
}
