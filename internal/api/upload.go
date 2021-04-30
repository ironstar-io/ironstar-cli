package api

import (
	"gitlab.com/ironstar-io/ironstar-cli/cmd/flags"
	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/console"
	"gitlab.com/ironstar-io/ironstar-cli/internal/system/fs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

// UploadPackage - Create a project tarball in tmp
func UploadPackage(creds types.Keylink, subHash, tarpath string, flg flags.Accumulator) (*RawResponse, error) {
	if flg.Tag != "" && flg.Branch != "" {
		return nil, errors.New("The fields 'branch' and 'tag' should not be specified at the same time.")
	}

	if flg.CustomPackage != "" {
		color.Red(`Warning! This command uploads the contents of this directory to your Ironstar Subscription. Only paths listed under the "exclude" settings in your .ironstar/config.yml file will be excluded.

This means that any database files, .env files, or other potentially sensitive content that you have in this repository that is not in the exclude list will be uploaded to the remote environment and possibly made publicly visible.

Please proceed with caution.
		`)
	} else {
		color.Red(`Warning! This command uploads the specified tarball to your Ironstar Subscription. Paths listed under the "exclude" settings in your .ironstar/config.yml are ignored when --custom-package is set.

This means that any database files, .env files, or other potentially sensitive content that you have in this tarball will be uploaded to the remote environment and possibly made publicly visible.

Please proceed with caution.
		`)
	}

	wo := console.SpinStart("Uploading package")

	req := &Stream{
		RunTokenRefresh: true,
		Credentials:     creds,
		Method:          "POST",
		FilePath:        tarpath,
		Path:            "/upload/subscription/" + subHash,
		Payload: map[string]string{
			"ref":      flg.Ref,
			"branch":   flg.Branch,
			"tag":      flg.Tag,
			"checksum": flg.Checksum,
		},
	}

	res, err := req.Send()

	if customPackage != "" {
		// Remove the tarball, regardless of the result. (not for custom packages)
		fs.Remove(tarpath)
	}

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
