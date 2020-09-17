package types

type PostBackupRequestParams struct {
	SubscriptionID string
	EnvironmentID  string
	Name           string
	Kind           string
	Components     []string
}
