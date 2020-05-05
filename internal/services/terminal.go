package services

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

// StdinPrompt ...
func StdinPrompt(prompt string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(text), nil
}

// StdinSecret ...
func StdinSecret(prompt string) (string, error) {
	fmt.Print(prompt)
	byteSecret, err := terminal.ReadPassword(0)
	if err != nil {
		return "", err
	}
	secret := string(byteSecret)

	return strings.TrimSpace(secret), nil
}

func GetCLIEmail(args []string) (string, error) {
	var email string
	if len(args) == 0 {
		input, err := StdinPrompt("Email: ")
		if err != nil {
			return "", err
		}
		email = input
	} else {
		email = args[0]
	}

	err := ValidateEmail(email)
	if err != nil {
		return "", err
	}

	return email, nil
}

func GetCLIPassword(passwordFlag string) (string, error) {
	var password string
	if passwordFlag == "" {
		input, err := StdinSecret("Password: ")
		if err != nil {
			return "", err
		}
		password = input

		fmt.Println()
	} else {
		color.Yellow("Warning: Supplying a password via the command line flag is potentially insecure")

		password = passwordFlag
	}

	return password, nil
}

func GetCLIProjectName() (string, error) {
	pname, err := StdinPrompt("Project Name: ")
	if err != nil {
		return "", err
	}

	return pname, nil
}

func GetCLIMFAPasscode() (string, error) {
	passcode, err := StdinSecret("MFA Passcode: ")
	if err != nil {
		return "", err
	}

	if len(passcode) != 6 {
		fmt.Println()
		color.Red("Ironstar API authentication failed!")
		fmt.Println()
		fmt.Println("MFA passcode length must be 6")

		return "", errors.New("Passcode length must be 6")
	}

	return passcode, nil
}

// GetDeployID
func GetDeployID(args []string, deployFlag string) (string, error) {
	if deployFlag != "" {
		return deployFlag, nil
	}

	if len(args) == 0 {
		di, err := StdinPrompt("Deployment ID: ")
		if err != nil {
			return "", errors.New("No deployment ID argument supplied")
		}

		return di, nil
	}

	return args[0], nil
}

// ConfirmationPrompt - The 'weighting' param should be one of [ "y", "n" ].
func ConfirmationPrompt(prompt, weighting string, autoAccept bool) bool {
	// Should be set by flag --yes (-y), allowing these prompts to be skipped
	if autoAccept {
		return true
	}

	response, err := StdinPrompt(prompt + weightedString(weighting))
	if err != nil {
		return false
	}

	var cutResponse string
	if response == "" {
		cutResponse = weighting
	} else {
		cutResponse = strings.ToLower(string([]rune(response)[0]))
	}

	if weighting == "n" {
		if cutResponse == "y" {
			return true
		}

		return false
	}

	if cutResponse == "n" {
		return false
	}

	return true
}

func weightedString(weighting string) string {
	if weighting == "y" {
		return " (Y/n): "
	}

	return " (y/N): "
}
