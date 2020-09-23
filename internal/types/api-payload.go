package types

type PostBackupRequestParams struct {
	SubscriptionID string
	EnvironmentID  string
	Name           string
	Kind           string
	LockTables     bool
	Components     []string
}

type DeleteBackupIterationParams struct {
	SubscriptionID string
	EnvironmentID  string
	Name           string
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
