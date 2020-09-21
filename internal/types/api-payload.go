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

type PostSyncRequestParams struct {
	SubscriptionID  string
	RestoreStrategy string
	SrcEnvironment  string
	DestEnvironment string
	Components      []string
}
