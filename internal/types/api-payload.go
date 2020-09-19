package types

type PostBackupRequestParams struct {
	SubscriptionID string
	EnvironmentID  string
	Name           string
	Kind           string
	Components     []string
}

type PostRestoreRequestParams struct {
	SubscriptionID string
	EnvironmentID  string
	Name           string
	Strategy       string
	Backup         string
	Components     []string
}
