package integration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
	. "github.com/paketo-buildpacks/occam/matchers"
)

func testSimpleApp(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect     = NewWithT(t).Expect
		Eventually = NewWithT(t).Eventually

		pack   occam.Pack
		docker occam.Docker
	)

	it.Before(func() {
		pack = occam.NewPack()
		docker = occam.NewDocker()

	})

	context("when building a simple app", func() {
		var (
			image     occam.Image
			container occam.Container

			name   string
			source string
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())

			source, err = occam.Source(filepath.Join("testdata", "passenger_sinatra"))
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		it("creates a working OCI image with a passenger start command", func() {
			var logs fmt.Stringer
			var err error

			image, logs, err = pack.WithNoColor().Build.
				WithBuildpacks(
					settings.Buildpacks.MRI.Online,
					settings.Buildpacks.Bundler.Online,
					settings.Buildpacks.BundleInstall.Online,
					settings.Buildpacks.Passenger.Online,
				).
				WithEnv(map[string]string{"BP_LOG_LEVEL": "DEBUG"}).
				WithPullPolicy("never").
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

			container, err = docker.Container.Run.
				WithEnv(map[string]string{"PORT": "9090"}).
				WithPublish("9090").
				WithPublishAll().
				Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())

			Eventually(container).Should(BeAvailable())
			Eventually(container).Should(Serve(ContainSubstring("Hello world!")).OnPort(9090))

			Expect(logs).To(ContainLines(
				MatchRegexp(fmt.Sprintf(`%s \d+\.\d+\.\d+`, settings.Buildpack.Name)),
				"  Executing build process",
				"  Getting the layer associated with curl:",
				fmt.Sprintf("    /layers/%s/curl", strings.ReplaceAll(settings.Buildpack.ID, "/", "_")),
			))
			Expect(logs).To(ContainLines(
				MatchRegexp(`    Installing cURL \d+\.\d+\.\d+`),
				fmt.Sprintf("    Installation path: /layers/%s/curl", strings.ReplaceAll(settings.Buildpack.ID, "/", "_")),
				MatchRegexp(`    Dependency URI\: https\:\/\/\S+\/curl\/curl_\d+\.\d+\.\d+_linux.*\.tgz`),
				MatchRegexp(`      Completed in ([0-9]*(\.[0-9]*)?[a-z]+)+`),
			))
			Expect(logs).To(ContainLines(
				"  Assigning launch processes:",
				"    web (default): bash -c bundle exec passenger start --port ${PORT:-3000}",
			))
		})

		context("validating SBOM", func() {
			var (
				sbomDir string
			)

			it.Before(func() {
				var err error
				sbomDir, err = os.MkdirTemp("", "sbom")
				Expect(err).NotTo(HaveOccurred())
				Expect(os.Chmod(sbomDir, os.ModePerm)).To(Succeed())
			})

			it.After(func() {
				Expect(os.RemoveAll(sbomDir)).To(Succeed())
			})

			it("writes SBOM files to the layer and label metadata", func() {
				var err error
				var logs fmt.Stringer

				image, logs, err = pack.WithNoColor().Build.
					WithPullPolicy("never").
					WithBuildpacks(
						settings.Buildpacks.MRI.Online,
						settings.Buildpacks.Bundler.Online,
						settings.Buildpacks.BundleInstall.Online,
						settings.Buildpacks.Passenger.Online,
					).
					WithEnv(map[string]string{
						"BP_LOG_LEVEL": "DEBUG",
					}).
					WithSBOMOutputDir(sbomDir).
					Execute(name, source)
				Expect(err).ToNot(HaveOccurred(), logs.String)

				container, err = docker.Container.Run.
					WithEnv(map[string]string{"PORT": "9090"}).
					WithPublish("9090").
					WithPublishAll().
					Execute(image.ID)
				Expect(err).NotTo(HaveOccurred())

				Eventually(container).Should(BeAvailable())
				Eventually(container).Should(Serve(ContainSubstring("Hello world!")).OnPort(9090))

				Expect(logs).To(ContainLines(
					fmt.Sprintf("  Generating SBOM for /layers/%s/curl", strings.ReplaceAll(settings.Buildpack.ID, "/", "_")),
					MatchRegexp(`      Completed in \d+(\.?\d+)*`),
				))
				Expect(logs).To(ContainLines(
					"  Writing SBOM in the following format(s):",
					"    application/vnd.cyclonedx+json",
					"    application/spdx+json",
					"    application/vnd.syft+json",
				))

				// check that all required SBOM files are present
				Expect(filepath.Join(sbomDir, "sbom", "launch", strings.ReplaceAll(settings.Buildpack.ID, "/", "_"), "curl", "sbom.cdx.json")).To(BeARegularFile())
				Expect(filepath.Join(sbomDir, "sbom", "launch", strings.ReplaceAll(settings.Buildpack.ID, "/", "_"), "curl", "sbom.spdx.json")).To(BeARegularFile())
				Expect(filepath.Join(sbomDir, "sbom", "launch", strings.ReplaceAll(settings.Buildpack.ID, "/", "_"), "curl", "sbom.syft.json")).To(BeARegularFile())

				// check an SBOM file to make sure it has an entry for curl
				contents, err := os.ReadFile(filepath.Join(sbomDir, "sbom", "launch", strings.ReplaceAll(settings.Buildpack.ID, "/", "_"), "curl", "sbom.cdx.json"))
				Expect(err).NotTo(HaveOccurred())
				Expect(string(contents)).To(ContainSubstring(`"name": "cURL"`))
			})
		})
	})
}
