package main

import (
	"os"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/cargo"
	"github.com/paketo-buildpacks/packit/v2/chronos"
	"github.com/paketo-buildpacks/packit/v2/postal"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/paketo-buildpacks/passenger"
)

func main() {
	packit.Run(
		passenger.Detect(passenger.NewGemfileParser()),
		passenger.Build(
			postal.NewService(cargo.NewTransport()),
			chronos.DefaultClock,
			scribe.NewEmitter(os.Stdout),
		),
	)
}
