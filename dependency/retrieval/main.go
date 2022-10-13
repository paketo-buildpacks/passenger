package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/joshuatcasey/collections"
	"github.com/joshuatcasey/libdependency/github"
	"github.com/joshuatcasey/libdependency/retrieve"
	"github.com/joshuatcasey/libdependency/upstream"
	"github.com/joshuatcasey/libdependency/versionology"
	"github.com/paketo-buildpacks/packit/v2/cargo"
	"golang.org/x/crypto/openpgp"
)

type StackAndTargetPair struct {
	stacks []string
	target string
}

var supportedStacks = []StackAndTargetPair{
	{stacks: []string{"io.buildpacks.stacks.jammy"}, target: "jammy"},
	{stacks: []string{"io.buildpacks.stacks.bionic"}, target: "bionic"},
}

func generateMetadata(versionFetcher versionology.VersionFetcher) ([]versionology.Dependency, error) {
	version := versionFetcher.Version().String()

	sourceURL := fmt.Sprintf("https://curl.se/download/curl-%s.tar.gz", version)
	ascURL := fmt.Sprintf("https://curl.se/download/curl-%s.tar.gz.asc", version)

	sourceSignature, err := getAsString(ascURL)
	if err != nil {
		return nil, err
	}

	curlTarballPath, err := downloadToFile(sourceURL)
	if err != nil {
		return nil, err
	}

	err = verifyASC(sourceSignature, curlTarballPath, danielHaxxKey)
	if err != nil {
		return nil, err
	}

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
		Version:        version,
	}

	return collections.TransformFuncWithError(supportedStacks, func(pair StackAndTargetPair) (versionology.Dependency, error) {
		configMetadataDependency.Stacks = pair.stacks
		return versionology.NewDependency(configMetadataDependency, pair.target)
	})
}

func verifyASC(signature, target, pgpKey string) error {
	file, err := os.Open(target)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	keyring, err := openpgp.ReadArmoredKeyRing(strings.NewReader(pgpKey))
	if err != nil {
		return err
	}

	signer, err := openpgp.CheckArmoredDetachedSignature(keyring, file, strings.NewReader(signature))
	if signer == nil {
		return errors.New("signature not accepted")
	}
	return err
}

func getAsString(url string) (string, error) {
	response, err := http.DefaultClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("could not get project metadata: %w", err)
	}

	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("could not read response: %w", err)
	}

	return string(body), nil
}

func downloadToFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to query url: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to query url %s with: status code %d", url, resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	tempDir, err := os.MkdirTemp("", "temp")
	if err != nil {
		return "", err
	}

	tempFilePath := filepath.Join(tempDir, filepath.Base(url))
	err = os.WriteFile(tempFilePath, body, os.ModePerm)
	if err != nil {
		return "", err
	}

	return tempFilePath, nil
}

func main() {
	getAllVersions := github.GetAllVersions(os.Getenv("GIT_TOKEN"), "curl", "curl")

	retrieve.NewMetadata("curl", getAllVersions, generateMetadata)
}
