package errs

import (
	"errors"
)

const (
	ShowCredentialsErrorMsg = "Unable to display credentials"

	DefaultCredentialsErrorMsg = "Unable to set default credentials"

	APILoginErrorMsg = "Ironstar API authentication failed"

	SetCredentialsErrorMsg = "Unable to set credentials"

	GetCredentialsErrorMsg = "Unable to get credentials"

	NoProjectFoundErrorMsg = "This command can only be run from a Tokaido project directory"

	NoSuitableCredsMsg = "There were no suitable credentials found for this project. Have you run `tok auth login`?"

	APISubLinkErrorMsg = "Ironstar API failed to link subscription"

	ProjectRootNotFoundError = "Unable to find a matching project"
)

var NoCredentialMatch = errors.New("There are no credentials available for the supplied email")

var NoSuitableCreds = errors.New(NoSuitableCredsMsg)
