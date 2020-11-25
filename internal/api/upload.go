package api

import (
	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/console"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/fs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

// UploadPackage - Create a project tarball in tmp
func UploadPackage(creds types.Keylink, subHash, tarpath, ref string) (*RawResponse, error) {
	color.Red(`Warning! This command uploads the contents of this Git repository to Ironstar and makes it available on the web. Only paths listed under the "exclude" settings in your .ironstar/config.yml file will be excluded.

This means that any database files, Drupal Public and Private Files, and other potentially sensitive content that you have in this repository that is not in the exclude list will be uploaded to the remote environment. 

Please proceed with caution. 
	`)

	wo := console.SpinStart("Uploading package file to Ironstar system")

	req := &Stream{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "POST",
		FilePath:         tarpath,
		Path:             "/upload/subscription/" + subHash,
		MapStringPayload: map[string]interface{}{},
		Ref:              ref,
	}

	res, err := req.Send()

	// Remove the tarball, regardless of the result.
	fs.Remove(tarpath)

	if err != nil {
		console.SpinPersist(wo, "⛔", "There was an error while uploading your package\n")
		return nil, errors.Wrap(err, errs.APISubListErrorMsg)
	}

	if res.StatusCode != 200 {
		console.SpinPersist(wo, "⛔", "There was an error while uploading your package\n")
		return nil, res.HandleFailure()
	}

	console.SpinPersist(wo, "💾", "Package upload completed successfully!\n")

	return res, nil
}
