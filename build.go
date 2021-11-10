package passenger

import (
	"path/filepath"
	"time"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/chronos"
	"github.com/paketo-buildpacks/packit/postal"
	"github.com/paketo-buildpacks/packit/scribe"
)

//go:generate faux --interface DependencyManager --output fakes/dependency_manager.go
type DependencyManager interface {
	Resolve(path, id, version, stack string) (postal.Dependency, error)
	Install(dependency postal.Dependency, cnbPath, layerPath string) error
}

func Build(dependencyManager DependencyManager, clock chronos.Clock, logger scribe.Logger) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		logger.Process("Executing build process")
		dependency, err := dependencyManager.Resolve(filepath.Join(context.CNBPath, "buildpack.toml"), "curl", "*", context.Stack)
		if err != nil {
			return packit.BuildResult{}, err
		}

		curlLayer, err := context.Layers.Get("curl")
		if err != nil {
			return packit.BuildResult{}, err
		}
		curlLayer.Launch = true

		logger.Subprocess("Installing %s %s", dependency.Name, dependency.Version)
		duration, err := clock.Measure(func() error {
			return dependencyManager.Install(dependency, context.CNBPath, curlLayer.Path)
		})
		if err != nil {
			return packit.BuildResult{}, err
		}
		logger.Action("Completed in %s", duration.Round(time.Millisecond))
		logger.Break()

		command := `bundle exec passenger start --port ${PORT:-3000}`
		logger.Process("Assigning launch processes")
		logger.Subprocess("web: %s", command)
		logger.Break()

		return packit.BuildResult{
			Layers: []packit.Layer{curlLayer},
			Launch: packit.LaunchMetadata{
				Processes: []packit.Process{
					{
						Type:    "web",
						Command: command,
						Default: true,
					},
				},
			},
		}, nil
	}
}
