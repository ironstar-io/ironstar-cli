package api

import (
	"fmt"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/system/console"
	"github.com/ironstar-io/ironstar-cli/internal/system/fs"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

// UploadPackage - Create a project tarball in tmp
func UploadPackage(creds types.Keylink, subHash, tarpath string, flg flags.Accumulator) (*RawResponse, error) {
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

	if flg.CustomPackage == "" {
		// Remove the tarball, regardless of the result. (not for custom packages)
		defer fs.Remove(tarpath)
	}

	retries := 2

	req := &Stream{
		RunTokenRefresh: true,
		Credentials:     creds,
		Method:          "POST",
		FilePath:        tarpath,
		URL:             fmt.Sprintf("%s/upload/subscription/%s", GetBaseUploadURL(), subHash),
		Payload: map[string]string{
			"ref":        flg.Ref,
			"branch":     flg.Branch,
			"tag":        flg.Tag,
			"checksum":   flg.Checksum,
			"commit_sha": flg.CommitSHA,
		},
	}

	res, err := retryHTTPWithExpBackoff(
		func() (*RawResponse, error) {
			wo := console.SpinStart("Uploading package")
			res, err := req.Send()
			if err != nil {
				console.SpinPersist(wo, "⛔", "There was an error while uploading the package\n")
				return nil, err
			}

			console.SpinPersist(wo, "💾", "Package upload completed successfully!\n")
			return res, nil
		}, retries)
	if err != nil {
		debugLogs(req.URL, retries, err)

		return nil, errors.Wrap(err, errs.UploadFailedErrorMsg)
	}

	if res.StatusCode != 200 {
		return nil, res.HandleFailure(flg.Output)
	}

	return res, nil
}
