package auth

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func IronstarSetActiveCredentials(args []string, flg flags.Accumulator) error {
	email, err := services.GetCLIEmail(args)
	if err != nil {
		return errors.Wrap(err, errs.SetCredentialsErrorMsg)
	}

	c, err := services.UpdateActiveCredentials(email)
	if err != nil {
		return errors.Wrap(err, errs.SetCredentialsErrorMsg)
	}

	if strings.ToLower(flg.Output) == "json" {
		utils.PrintInterfaceAsJSON(c)
		return nil
	}

	fmt.Println()
	color.Green("Ironstar API user now active!")
	fmt.Println()
	color.Green("User: ")
	fmt.Println(email)

	if c.Expiry.IsZero() {
		// Should always return expiry, but check anyway for safety
		return nil
	}

	fmt.Println()
	color.Green("Expiry: ")

	expDiff := strconv.Itoa(int(math.RoundToEven(c.Expiry.Sub(time.Now().UTC()).Hours())))
	fmt.Println(c.Expiry.String() + " (" + expDiff + " hours)")

	return nil
}
