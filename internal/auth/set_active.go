package auth

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"gitlab.com/ironstar-io/ironstar-cli/internal/errs"
	"gitlab.com/ironstar-io/ironstar-cli/internal/services"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func IronstarSetActiveCredentials(args []string) error {
	email, err := services.GetCLIEmail(args)
	if err != nil {
		return errors.Wrap(err, errs.SetCredentialsErrorMsg)
	}

	c, err := services.UpdateActiveCredentials(email)
	if err != nil {
		return errors.Wrap(err, errs.SetCredentialsErrorMsg)
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

	expDiff := strconv.Itoa(int(math.RoundToEven(c.Expiry.Sub(time.Now().UTC()).Hours() / 24)))
	fmt.Println(c.Expiry.String() + " (" + expDiff + " days)")

	return nil
}
