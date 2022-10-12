package main

import (
	"fmt"
	"os"

	"github.com/joshuatcasey/libdependency/github"
	"github.com/joshuatcasey/libdependency/retrieve"
	"github.com/joshuatcasey/libdependency/upstream"
	"github.com/joshuatcasey/libdependency/versionology"
	"github.com/paketo-buildpacks/packit/v2/cargo"
)

func generateMetadata(versionFetcher versionology.VersionFetcher) ([]versionology.Dependency, error) {
	version := versionFetcher.Version().String()

	sourceURL := fmt.Sprintf("https://curl.se/download/curl-%s.tar.gz", version)
	sourceSHA256, err := upstream.GetSHA256OfRemoteFile(sourceURL)
	if err != nil {
		return nil, err
	}

	configMetadataDependency := cargo.ConfigMetadataDependency{
		CPE:            fmt.Sprintf("cpe:2.3:a:haxx:curl:%s:*:*:*:*:*:*:*", version),
		ID:             "curl",
		Name:           "cURL",
		Licenses:       retrieve.LookupLicenses(sourceURL, upstream.DefaultDecompress),
		PURL:           retrieve.GeneratePURL("curl", version, sourceSHA256, sourceURL),
		Source:         sourceURL,
		SourceChecksum: fmt.Sprintf("sha256:%s", sourceSHA256),
		Stacks:         []string{"io.buildpacks.stacks.bionic"},
		Version:        version,
	}

	return versionology.NewDependencyArray(configMetadataDependency, "bionic")
}

func main() {
	getAllVersions := github.GetAllVersions(os.Getenv("GIT_TOKEN"), "curl", "curl")

	retrieve.NewMetadata("curl", getAllVersions, generateMetadata)
}
