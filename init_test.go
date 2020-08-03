package passenger_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitBundleInstall(t *testing.T) {
	suite := spec.New("passenger", spec.Report(report.Terminal{}), spec.Parallel())
	suite("Detect", testDetect)
	suite.Run(t)
}
