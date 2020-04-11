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

	ProjectRootNotFoundError = "Unable to find a matching project"

	APISubListErrorMsg = "Failed to retrieve subscriptions"

	APISubLinkErrorMsg = "Failed to link subscription"

	APIGetSubscriptionErrorMsg = "Failed to get subscription"

	UnexpectedErrorMsg = "An unexpected error occurred"
)

var NoCredentialMatch = errors.New("There are no credentials available for the supplied email")

var UnexpectedError = errors.New(UnexpectedErrorMsg)

var NoSuitableCreds = errors.New(NoSuitableCredsMsg)
