package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/system/console"
	"github.com/ironstar-io/ironstar-cli/internal/system/fs"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/pkg/errors"
)

// UploadPackage - Create a project tarball in tmp
func UploadPackage(creds types.Keylink, subHash, tarpath string, flg flags.Accumulator) (*RawResponse, error) {
	if flg.CustomPackage != "" {
		// The normal-package upload notice (with the exclude source) is printed by
		// services.CreateProjectTar; --custom-package skips that path, so cover it here.
		fmt.Printf("This will upload the custom package %s to Ironstar. Exclude rules don't apply to --custom-package.\n", flg.CustomPackage)
		fmt.Println("Sensitive files in the tarball (e.g. .env files or database dumps) will be included — review before continuing.")
	}

	if flg.CustomPackage == "" {
		// Remove the tarball, regardless of the result. (not for custom packages)
		defer fs.Remove(tarpath)
	}

	retries := 2

	req := &Stream{
		RunTokenRefresh: true,
		Credentials:     creds,
		Method:          http.MethodPost,
		FilePath:        tarpath,
		URL:             GetUploadURL(subHash),
		Payload: map[string]string{
			"ref":        utils.TruncateString(flg.Ref, 255),
			"branch":     utils.TruncateString(flg.Branch, 255),
			"tag":        utils.TruncateString(flg.Tag, 255),
			"checksum":   utils.TruncateString(flg.Checksum, 255),
			"commit_sha": utils.TruncateString(flg.CommitSHA, 255),
		},
	}

	res, err := retryHTTPWithExpBackoff(
		func() (*RawResponse, error) {
			wo := console.SpinStart("Uploading package")
			start := time.Now()
			res, err := req.Send()
			elapsed := time.Since(start).Round(time.Millisecond)
			if err != nil {
				console.SpinPersist(wo, "⛔", fmt.Sprintf("Package upload failed after %s with no HTTP response: %s\n", elapsed, err))
				return nil, err
			}
			if res.StatusCode < 199 || res.StatusCode > 299 {
				console.SpinPersist(wo, "⛔", fmt.Sprintf("Package upload failed after %s with response code %d\n", elapsed, res.StatusCode))
				return nil, fmt.Errorf("an error occurred with status %d: %s", res.StatusCode, res.Body)
			}

			console.SpinPersist(wo, "💾", "Package upload completed successfully!\n")
			return res, nil
		}, retries)
	if err != nil {
		debugLogs(req.URL, retries, err)

		return nil, errors.Wrap(err, errs.UploadFailedErrorMsg)
	}

	if res.StatusCode < 199 || res.StatusCode > 299 {
		return nil, res.HandleFailure(flg.Output)
	}

	return res, nil
}
