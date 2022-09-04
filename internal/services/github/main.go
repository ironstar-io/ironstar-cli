package github

import (
	"github.com/ironstar-io/ironstar-cli/internal/api"
)

// ReleaseBody - Used properties from the GH API GET call
type ReleaseBody struct {
	TagName    string `json:"tag_name"`
	HTMLURL    string `json:"html_url"`
	Prerelease bool   `json:"prerelease"`
	Draft      bool   `json:"draft"`
	Assets     []struct {
		Name string `json:"name"`
		URL  string `json:"browser_download_url"`
	}
}

// GetRelease - Get an Ironstar CLI release
func GetRelease(version string) (*api.RawResponse, error) {
	req := api.Request{
		Method:           "GET",
		URL:              "https://api.github.com/repos/ironstar-io/ironstar-cli/releases" + version,
		MapStringPayload: map[string]interface{}{},
	}
	res, err := req.HTTPSend()
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, res.HandleExternalFailure()
	}

	return res, nil
}

// GetLatestRelease - Get latest Ironstar CLI release
func GetLatestRelease() (*api.RawResponse, error) {
	return GetRelease("/latest")
}

// GetAllReleases - Get all Ironstar CLI releases
func GetAllReleases() (*api.RawResponse, error) {
	return GetRelease("")
}
