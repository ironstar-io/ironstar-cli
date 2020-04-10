package subscription

import (
	"path/filepath"

	"gitlab.com/ironstar-io/ironstar-cli/internal/api"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/fs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func Link(args []string) error {
	proj, err := services.GetProjectData()
	if err != nil {
		return err
	}

	user, err := services.ResolveUserCredentials(proj)
	if err != nil {
		return err
	}

	req := &api.Request{
		AuthToken:        user.AuthToken,
		Method:           "GET",
		Path:             "/user/subscriptions",
		MapStringPayload: map[string]string{},
	}

	res, err := req.Send()
	if err != nil {
		return errors.Wrap(err, APISubListErrorMsg)
	}

	if res.StatusCode != 200 {
		return res.HandleFailure()
	}

	var sub types.Subscription
	err = yaml.Unmarshal(res.Body, &sub)
	if err != nil {
		return err
	}

	return LinkSubscriptionToProject(proj, sub)
}

func LinkSubscriptionToProject(proj types.ProjectConfig, sub types.Subscription) error {
	projConf, err := services.ReadInProjectConfig(proj.Path)
	if err != nil {
		return errors.Wrap(err, APISubListErrorMsg)
	}

	projConf.Subscription = sub

	newMarhsalled, err := yaml.Marshal(projConf)
	if err != nil {
		return err
	}

	py := filepath.Join(proj.Path, ".ironstar", "global.yml")
	fs.Replace(py, newMarhsalled)

	return nil
}
