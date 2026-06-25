package pkg

import (
	"fmt"
	"strings"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func Create(args []string, flg flags.Accumulator) error {
	// A dry run previews what would be packaged and uploads nothing, so it needs
	// no credentials or subscription context — keep it a purely local operation.
	if flg.DryRun {
		indexPath, total, count, err := services.WritePackageIndex(flg)
		if err != nil {
			return err
		}

		if strings.ToLower(flg.Output) == "json" {
			utils.PrintInterfaceAsJSON(map[string]any{
				"dryRun":            true,
				"indexFile":         indexPath,
				"fileCount":         count,
				"uncompressedBytes": total,
			})
			return nil
		}

		fmt.Println()
		fmt.Printf("Dry run — no upload performed.\nWould package %d files, %s (uncompressed).\nIndex written to %s\n", count, humanize.IBytes(uint64(total)), indexPath)
		return nil
	}

	creds, err := services.ResolveUserCredentials(flg.Login)
	if err != nil {
		return err
	}

	sub, err := api.GetSubscriptionContext(creds, flg)
	if err != nil {
		return err
	}

	if sub.Alias == "" {
		return errs.ErrNoSubLink
	}

	utils.PrintCommandContext(flg.Output, creds.Login, sub.Alias, sub.HashedID)

	if flg.Tag != "" && flg.Branch != "" {
		return errors.New("The fields 'branch' and 'tag' should not be specified at the same time.")
	}

	tarpath, err := services.CreateProjectTar(flg)
	if err != nil {
		return err
	}

	res, err := api.UploadPackage(creds, sub.HashedID, tarpath, flg)
	if err != nil {
		return err
	}

	var ur types.UploadResponse
	err = yaml.Unmarshal(res.Body, &ur)
	if err != nil {
		return err
	}

	if strings.ToLower(flg.Output) == "json" {
		utils.PrintInterfaceAsJSON(ur)
		return nil
	}

	// Ensure both BuildID and HashedId are populated
	if ur.BuildID == "" && ur.HashedId != "" {
		ur.BuildID = ur.HashedId
	} else if ur.HashedId == "" && ur.BuildID != "" {
		ur.HashedId = ur.BuildID
	}

	// Ensure both BuildName and Name are populated
	if ur.BuildName == "" && ur.Name != "" {
		ur.BuildName = ur.Name
	} else if ur.Name == "" && ur.BuildName != "" {
		ur.Name = ur.BuildName
	}

	fmt.Println("PACKAGE ID: " + ur.BuildID)
	fmt.Println("PACKAGE NAME: " + ur.BuildName)

	return nil
}
