package passenger

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/chronos"
	"github.com/paketo-buildpacks/packit/v2/postal"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

//go:generate faux --interface DependencyManager --output fakes/dependency_manager.go
//go:generate faux --interface PassengerfileConfigParser --output fakes/passengerfile_parser.go

type DependencyManager interface {
	Resolve(path, id, version, stack string) (postal.Dependency, error)
	Deliver(dependency postal.Dependency, cnbPath, layerPath, platformPath string) error
}

type PassengerfileConfigParser interface {
	ParsePort(path string, defaultPort int) (int, error)
}

func Build(dependencyManager DependencyManager, passengerfileParser PassengerfileConfigParser, clock chronos.Clock, logger scribe.Emitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		logger.Process("Executing build process")
		dependency, err := dependencyManager.Resolve(filepath.Join(context.CNBPath, "buildpack.toml"), "curl", "*", context.Stack)
		if err != nil {
			return packit.BuildResult{}, err
		}

		logger.Debug.Process("Getting the layer associated with curl:")
		curlLayer, err := context.Layers.Get("curl")
		if err != nil {
			return packit.BuildResult{}, err
		}
		logger.Debug.Subprocess(curlLayer.Path)
		logger.Debug.Break()
		curlLayer.Launch = true

		logger.Subprocess("Installing %s %s", dependency.Name, dependency.Version)
		duration, err := clock.Measure(func() error {
			logger.Debug.Subprocess("Installation path: %s", curlLayer.Path)
			logger.Debug.Subprocess("Dependency URI: %s", dependency.URI)
			return dependencyManager.Deliver(dependency, context.CNBPath, curlLayer.Path, context.Platform.Path)
		})
		if err != nil {
			return packit.BuildResult{}, err
		}
		logger.Action("Completed in %s", duration.Round(time.Millisecond))
		logger.Break()

		passengerfilePath := filepath.Join(context.WorkingDir, "Passengerfile.json")

		defaultPort := 3000
		port, err := passengerfileParser.ParsePort(passengerfilePath, defaultPort)
		if err != nil {
			return packit.BuildResult{}, err
		}

		args := fmt.Sprintf(`bundle exec passenger start --port ${PORT:-%d}`, port)
		processes := []packit.Process{
			{
				Type:    "web",
				Command: "bash",
				Args:    []string{"-c", args},
				Default: true,
				Direct:  true,
			},
		}
		logger.LaunchProcesses(processes)

		return packit.BuildResult{
			Layers: []packit.Layer{curlLayer},
			Launch: packit.LaunchMetadata{
				Processes: processes,
			},
		}, nil
	}
}
