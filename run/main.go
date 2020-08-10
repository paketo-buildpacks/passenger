package main

import (
	"os"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/cargo"
	"github.com/paketo-buildpacks/packit/chronos"
	"github.com/paketo-buildpacks/packit/postal"
	"github.com/paketo-buildpacks/packit/scribe"
	"github.com/paketo-community/passenger"
)

func main() {
	packit.Run(
		passenger.Detect(passenger.NewGemfileParser()),
		passenger.Build(
			postal.NewService(cargo.NewTransport()),
			chronos.DefaultClock,
			scribe.NewLogger(os.Stdout),
		),
	)
}
