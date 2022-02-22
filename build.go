package passenger

import (
	"path/filepath"
	"time"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/chronos"
	"github.com/paketo-buildpacks/packit/v2/postal"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

//go:generate faux --interface DependencyManager --output fakes/dependency_manager.go
type DependencyManager interface {
	Resolve(path, id, version, stack string) (postal.Dependency, error)
	Deliver(dependency postal.Dependency, cnbPath, layerPath, platformPath string) error
}

func Build(dependencyManager DependencyManager, clock chronos.Clock, logger scribe.Emitter) packit.BuildFunc {
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

		args := `bundle exec passenger start --port ${PORT:-3000}`
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
