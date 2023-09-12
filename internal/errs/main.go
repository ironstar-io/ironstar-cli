package errs

import (
	"errors"
)

const (
	ShowCredentialsErrorMsg = "Unable to display credentials"

	DefaultCredentialsErrorMsg = "Unable to set default credentials"

	APILoginErrorMsg = "Ironstar API authentication failed"

	APILogoutErrorMsg = "Ironstar API logout failed"

	APIMFAEnrolErrorMsg = "Ironstar API MFA enrol failed"

	APIMFAVerifyErrorMsg = "Ironstar API MFA verification failed"

	SetCredentialsErrorMsg = "Unable to set credentials"

	GetCredentialsErrorMsg = "Unable to get credentials"

	NoProjectFoundErrorMsg = "This command can only be run from an Ironstar project directory"

	NoSuitableCredsMsg = "There were no suitable credentials found for this project. Have you run `iron login`?"

	ProjectRootNotFoundError = "Unable to find a matching project"

	APISubListErrorMsg = "Failed to retrieve subscriptions"

	APISubLinkErrorMsg = "Failed to link subscription"

	APIGetSubscriptionErrorMsg = "Failed to get subscription"

	APIGetBackupErrorMsg = "Failed to get backup"

	APIQueryLogsErrorMsg = "Failed to retrieve logs"

	APIPostBackupErrorMsg = "Failed to create backup"

	APIDeleteBackupErrorMsg = "Failed to delete backup"

	APIGetRestoreErrorMsg = "Failed to get restore"

	APIPostRestoreErrorMsg = "Failed to create restore"

	APIGetSyncErrorMsg = "Failed to get sync"

	APIPostSyncErrorMsg = "Failed to create sync"

	APIGetEnvironmentErrorMsg = "Failed to get environment"

	APIGetEnvironmentVariablesErrorMsg = "Failed to get environment variables"

	APIGetRemoteCommandsErrorMsg = "Failed to get remote commands"

	APIPostRemoteCommandsErrorMsg = "Failed to create remote command"

	APIPostEnvironmentVariableErrorMsg = "Failed to add new environment variable"

	APIPostCacheInvalidationErrorMsg = "Failed to create cache invalidation"

	APIGetCacheInvalidationErrorMsg = "Failed to retrieve cache invalidations"

	APIGetAntivirusScanErrorMsg = "Failed to retrieve antivirus scans"

	APIDeleteEnvironmentVariableErrorMsg = "Failed to delete environment variable"

	APIUpdateEnvironmentErrorMsg = "Failed to update environment"

	NoSubscriptionLinkErrorMsg = "No Ironstar subscription has been linked to this project. Have you run `iron subscription link [subscription-name]`"

	NoEnvironmentFlagSupplied = "An environment flag must be supplied for this command with `--env=[env-name]`"

	UnexpectedErrorMsg = "An unexpected error occurred"

	UploadFailedErrorMsg = "An unexpected error occurred during upload"

	IronstarAPIConnectionErrorMsg = "Unable to connect to the Ironstar API. Please try again later"
)

var ErrNoCredentialMatch = errors.New("There are no credentials available for the supplied email")

var ErrUnexpectedError = errors.New(UnexpectedErrorMsg)

var ErrNoSuitableCreds = errors.New(NoSuitableCredsMsg)

var ErrNoSubLink = errors.New("No Ironstar subscription has been linked to this project. Have you run `iron subscription link [subscription-name]`")
