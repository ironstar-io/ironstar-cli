package auth

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/api"
	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/services"
	"github.com/ironstar-io/ironstar-cli/internal/system/fs"
	"github.com/ironstar-io/ironstar-cli/internal/system/utils"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func MFAEnable(flg flags.Accumulator, creds types.Keylink, email string) (*types.AuthResponseBody, error) {
	if creds == (types.Keylink{}) {
		cr, err := services.ResolveUserCredentials(flg.Login)
		if err != nil {
			return nil, err
		}
		creds = cr
	}

	req := &api.Request{
		RunTokenRefresh:  true,
		Credentials:      creds,
		Method:           "POST",
		Path:             "/auth/mfa/enrol",
		MapStringPayload: map[string]interface{}{},
	}

	res, err := req.NankaiSend()
	if err != nil {
		return nil, errors.Wrap(err, errs.APISubListErrorMsg)
	}

	if res.StatusCode != 200 {
		return nil, res.HandleFailure()
	}

	var m types.MFAEnrolResponse
	err = json.Unmarshal(res.Body, &m)
	if err != nil {
		return nil, err
	}
	if m.QRCode == "" {
		return nil, errors.New(errs.APIMFAEnrolErrorMsg)
	}

	// Display QR token
	color.Green("Opening the QR code for display in your browser")
	fmt.Println()

	p, err := generateQRHTMLPage(creds.Login, m.QRCode)
	if err != nil {
		return nil, err
	}

	output, err := utils.OpenSite(p)
	if err != nil {
		color.Yellow("Unable to automatically open a browser in order to display your MFA QR code")
		fmt.Println()
		color.Yellow(output)
		color.Yellow(err.Error())
		fmt.Println()
		color.Yellow("You can view your QR code manually by pasting the following file path in your browser address bar")
		fmt.Println()
		fmt.Println(p)
		fmt.Println()
		color.Yellow("Alternatively you can enter this code into your Authenticator application manually")
		fmt.Println()
		fmt.Println(m.Secret)
	}

	fmt.Println()
	color.Green("Once registered in your preferred Authenticator application, we'll need to verify a supplied passcode")

	// Request first passcode for verification
	c, err := validateMFAPasscodeWithRetries(flg, m.IDToken, email)
	if err != nil {
		color.Yellow("Please try running `iron auth mfa enable` again or reach out to Ironstar Support for help`")
		fmt.Println()
		return nil, err
	}

	err = services.UpdateCredentialsFile(types.Keylink{
		Login:     creds.Login,
		AuthToken: c.IDToken,
		Expiry:    c.Expiry,
	})
	if err != nil {
		return nil, errors.Wrap(err, errs.APIMFAVerifyErrorMsg)
	}

	fmt.Println()
	fmt.Println()
	color.Green("MFA Recovery Code:")
	fmt.Println(m.RecoveryCode)
	fmt.Println()
	color.Green("Successfully enabled MFA for this account!")

	return c, nil
}

func generateQRHTMLPage(login, qrImg string) (string, error) {
	hp := "/tmp/ironstar"
	fs.Mkdir(hp)

	tmpl := buildQRHTMLTemplate(qrImg)

	hf := filepath.Join(hp, "qr.html")

	fs.Remove(hf)

	err := fs.TouchByteArray(hf, []byte(tmpl), 0400)
	if err != nil {
		return "", err
	}

	return hf, nil
}

func buildQRHTMLTemplate(qrImg string) string {
	return `<!DOCTYPE html>
<html lang="en">
	<head>
	<meta charset="utf-8" />
	<title>Ironstar CLI - MFA QR Code</title>
	<style>
		body {
		font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
			Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji",
			"Segoe UI Symbol";
		text-align: center;
		}
	</style>
	</head>
	<body>
	<h1>Ironstar CLI QR Code</h1>
	<p>Scan the following code in your preferred Authenticator application</p>
	<img
		src="data:image/png;base64,` + qrImg + `"
	/>
	<p>
		Once completed you'll need to verify a passcode<br />
		back in your terminal to enable MFA
	</p>
	</body>
</html>`
}
