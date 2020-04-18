package api

import (
	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/pkg/errors"
)

// UploadPackage - Create a project tarball in tmp
func UploadPackage(creds types.Keylink, tarpath, subHash string) (*RawResponse, error) {
	req := &Stream{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "POST",
		FilePath:         tarpath,
		Path:             "/upload/subscription/" + subHash,
		MapStringPayload: map[string]string{},
	}

	res, err := req.Send()
	if err != nil {
		return nil, errors.Wrap(err, errs.APISubListErrorMsg)
	}

	if res.StatusCode != 200 {
		return nil, res.HandleFailure()
	}

	return res, nil
}
