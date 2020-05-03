package version

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"gitlab.com/ironstar-io/ironstar-cli/internal/system/console"

	"github.com/blang/semver"
)

// TODO
// Latest is the latest version response from release.tokaido.io/latest
type Latest struct {
	Version string `json:"version"`
	URL     string `json:"url"`
}

// Check - checks if the current ironstar-cli version is the latest available
func Check() error {
	// c := conf.GetConfig()
	info := Get()
	re := regexp.MustCompile("[^-]*")
	match := re.FindStringSubmatch(info.Version)
	current, _ := semver.Make(match[0])

	latestVersion, _, err := getLatestVersion()
	if err != nil {
		return err
	}
	latestSemver, _ := semver.Make(latestVersion)
	if err != nil {
		return err
	}

	if latestSemver.GT(current) {
		console.Println("\n👵🏻  You're running an old version of the Ironstar CLI. Please consider upgrading to Ironstar CLI "+latestVersion+"   ", "")
		console.Println("    You can upgrade easily by running 'iron upgrade'", "")
	}

	return nil
}

func getLatestVersion() (ver, url string, err error) {
	// TODO
	req, err := http.NewRequest(http.MethodGet, "https://api.tokaido.io/v1/release/latest", nil)
	if err != nil {
		return "", "", err
	}

	client := http.Client{
		Timeout: time.Second * 3,
	}
	res, err := client.Do(req)
	if err != nil {
		return "", "", err
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return "", "", err
	}

	latest := Latest{}
	err = json.Unmarshal(body, &latest)
	if err != nil {
		return "", "", err
	}

	return latest.Version, latest.URL, nil

}
