package main

import (
	"os"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/scribe"
	"github.com/paketo-community/passenger"
)

func main() {
	parser := passenger.NewGemfileParser()
	logger := scribe.NewLogger(os.Stdout)
	packit.Run(
		passenger.Detect(parser),
		passenger.Build(logger),
	)
}
