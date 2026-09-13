package types

type CacheInvalidationKind string
type CacheInvalidationType string

const (
	CacheInvalidationKindEnvironment CacheInvalidationKind = "environment"
	CacheInvalidationKindURL         CacheInvalidationKind = "url"
	CacheInvalidationTypeHard        CacheInvalidationType = "hard"
	CacheInvalidationTypeSoft        CacheInvalidationType = "soft"
)

type PostCacheInvalidationRequestParams struct {
	Kind             CacheInvalidationKind `json:"kind"`
	InvalidationType CacheInvalidationType `json:"invalidation_type"`
	URL              string                `json:"url,omitempty"`
}

type PostBackupRequestParams struct {
	SubscriptionID string
	EnvironmentID  string
	Name           string
	Kind           string
	LockTables     bool
	Components     []string
}

type DeleteBackupParams struct {
	SubscriptionID string
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

type PutNewRelicParams struct {
	LicenseKey  string `json:"licenseKey"`
	AppID       string `json:"appID"`
	AppName     string `json:"appName"`
	APIKeyValue string `json:"apiKeyValue"`
	APIKeyType  string `json:"apiKeyType,omitempty"`
}

type PostSyncRequestParams struct {
	SubscriptionID  string
	RestoreStrategy string
	SrcEnvironment  string
	DestEnvironment string
	Components      []string
}
