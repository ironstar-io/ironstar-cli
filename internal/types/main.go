package types

import (
	"reflect"
	"time"
)

type Keylink struct {
	Login     string    `json:"login"`
	AuthToken string    `json:"auth_token"`
	Expiry    time.Time `json:"expiry"`
}

type Credentials struct {
	Active   string    `json:"active"`
	Keychain []Keylink `json:"keychain"`
}

type AuthResponseBody struct {
	IDToken          string    `json:"id_token"`
	RedirectEndpoint string    `json:"redirect_endpoint"`
	MFAStatus        string    `json:"mfa_status"` // If this is set, user is MFA registered
	Expiry           time.Time `json:"expiry"`
}

type ProjectConfig struct {
	Version      string        `yaml:"version,omitempty"`
	Subscription Subscription  `yaml:"subscription,omitempty"`
	Package      PackageConfig `yaml:"package,omitempty"`
}

type PackageConfig struct {
	Exclude []string `yaml:"exclude,omitempty"`
}

type MFAEnrolResponse struct {
	QRCode       string    `json:"qr_code,omitempty"`
	IDToken      string    `json:"id_token"`
	Secret       string    `json:"secret"`
	RecoveryCode string    `json:"recovery_code"`
	Expiry       time.Time `json:"expiry"`
}

// Subscription
type Subscription struct {
	HashedID        string `json:"subscription_id,omitempty" yaml:"subscription_id,omitempty" example:"98hreHs"`
	Alias           string `json:"alias,omitempty" yaml:"alias,omitempty" example:"umami-food-blog"`
	Ref             string `json:"ref,omitempty" yaml:"ref,omitempty" example:"au1999"`
	ApplicationType string `json:"application_type,omitempty" yaml:"application_type,omitempty" example:"drupal"`
}

type SubscriptionEnvironment struct {
	Subscription
	Environment
}

type Role struct {
	Name        string   `json:"name,omitempty" yaml:"name,omitempty"`
	Permissions []string `json:"permissions,omitempty"  yaml:"permissions,omitempty"`
}

// SubscriptionAccessResponse
type UserAccessResponse struct {
	Role         Role         `json:"role,omitempty" yaml:"role,omitempty"`
	Subscription Subscription `json:"subscription,omitempty" yaml:"subscription,omitempty"`
	// Environment  Environment  `yaml:"environment,omitempty"`
}

// Build
type BuildFlat struct {
	HashedID  string    `json:"build_id,omitempty" yaml:"build_id,omitempty"`
	Name      string    `json:"name,omitempty" yaml:"name,omitempty"`
	Status    string    `json:"status,omitempty" yaml:"status,omitempty"`
	CreatedBy string    `json:"created_by,omitempty" yaml:"created_by,omitempty"`
	RunningIn string    `json:"running_in,omitempty" yaml:"running_in,omitempty"`
	Branch    string    `json:"branch,omitempty" yaml:"branch,omitempty"`
	Tag       string    `json:"tag,omitempty" yaml:"tag,omitempty"`
	Checksum  string    `json:"checksum,omitempty" yaml:"checksum,omitempty"`
	CommitSHA string    `json:"commit_sha,omitempty" yaml:"commit_sha,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty"`
}

// Build
type Build struct {
	HashedID   string       `json:"build_id,omitempty" yaml:"build_id,omitempty"`
	Name       string       `json:"name,omitempty" yaml:"name,omitempty"`
	Ref        string       `json:"ref,omitempty" yaml:"ref,omitempty"`
	Status     string       `json:"status,omitempty" yaml:"status,omitempty"`
	Branch     string       `json:"branch,omitempty" yaml:"branch,omitempty"`
	Tag        string       `json:"tag,omitempty" yaml:"tag,omitempty"`
	Checksum   string       `json:"checksum,omitempty" yaml:"checksum,omitempty"`
	CommitSHA  string       `json:"commit_sha,omitempty" yaml:"commit_sha,omitempty"`
	Notes      string       `json:"notes,omitempty" yaml:"notes,omitempty"`
	CreatedBy  string       `json:"created_by,omitempty" yaml:"created_by,omitempty"`
	RunningIn  string       `json:"running_in,omitempty" yaml:"running_in,omitempty"`
	CreatedAt  time.Time    `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	Deployment []Deployment `json:"deployments,omitempty" yaml:"deployments,omitempty"`
}

func (s Deployment) IsStructureEmpty() bool {
	return reflect.DeepEqual(s, Deployment{})
}

// Deployment
type Deployment struct {
	HashedID       string                     `json:"deployment_id,omitempty" yaml:"deployment_id,omitempty"`
	Name           string                     `json:"name,omitempty" yaml:"name,omitempty"`
	AppStatus      string                     `json:"app_status,omitempty" yaml:"app_status,omitempty"`
	AdminSvcStatus string                     `json:"admin_svc_status,omitempty" yaml:"admin_svc_status,omitempty"`
	BuildID        string                     `json:"build_token,omitempty" yaml:"build_token,omitempty"`
	CreatedAt      time.Time                  `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	Environment    Environment                `json:"environment,omitempty" yaml:"environment,omitempty"`
	Build          BuildFlat                  `json:"build,omitempty" yaml:"build,omitempty"`
	Status         DeploymentStatus           `json:"status,omitempty" yaml:"status,omitempty"`
	Lifecycle      []DeploymentLifecycleEvent `json:"lifecycle,omitempty" yaml:"lifecycle,omitempty"`
}

type DeploymentLifecycleEvent struct {
	Stage   string    `json:"stage,omitempty" example:"Requested" extensions:"x-description=The name of the lifecycle stage"`
	Status  string    `json:"status,omitempty" example:"Successful" extensions:"x-description=The status of the lifecycle change"`
	Command string    `json:"command,omitempty" example:"drush cr" extensions:"x-description=The command in the lifecycle change"`
	Host    string    `json:"host,omitempty" example:"host" extensions:"x-description=The host name of the deployment"`
	Enter   time.Time `json:"enter,omitempty" example:"2023-01-02T15:04:05Z07:00" format:"RFC3339" extensions:"x-description=The date in which the deployment entered this state"`
	Exit    time.Time `json:"exit,omitempty" example:"2023-01-02T15:04:05Z07:00" format:"RFC3339" extensions:"x-description=The date in which the deployment exited this state"`
}

type DeploymentStatus struct {
	Lifecycle      string    `json:"lifecycle,omitempty" yaml:"lifecycle,omitempty"`
	ActiveDate     time.Time `json:"activeDate,omitempty" yaml:"activeDate,omitempty"`
	RetiredDate    time.Time `json:"retiredDate,omitempty" yaml:"retiredDate,omitempty"`
	AppDeployments int       `json:"appDeployments,omitempty" yaml:"appDeployments,omitempty"`
}

// UserMinimal
type UserMinimal struct {
	FirstName   string `json:"first_name,omitempty" yaml:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty" yaml:"last_name,omitempty"`
	DisplayName string `json:"display_name,omitempty" yaml:"display_name,omitempty"`
	Email       string `json:"email,omitempty" yaml:"email,omitempty"`
}

// DeploymentActivityResponse
type DeploymentActivityResponse struct {
	Message   string    `json:"message,omitempty" yaml:"message,omitempty"`
	Flag      string    `json:"flag,omitempty" yaml:"flag,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty"`
}

type Environment struct {
	HashedID          string `json:"environment_id,omitempty" yaml:"environment_id,omitempty"`
	Class             string `json:"class,omitempty" yaml:"class,omitempty"`
	DNSName           string `json:"dns_name,omitempty" yaml:"dns_name,omitempty"`
	Name              string `json:"name,omitempty" yaml:"name,omitempty"`
	Trigger           string `json:"trigger,omitempty" yaml:"trigger,omitempty"`
	UpdateStrategy    string `json:"update_strategy,omitempty" yaml:"update_strategy,omitempty"`
	RestorePermission string `json:"restore_permission,omitempty" yaml:"restore_permission,omitempty"`
	LogRetention      int64  `json:"log_retention,omitempty" yaml:"log_retention,omitempty"`
}

type EnvVars struct {
	Key       string    `json:"key,omitempty" yaml:"key,omitempty"`
	Value     string    `json:"value,omitempty" yaml:"value,omitempty"`
	VarType   string    `json:"var_type,omitempty" yaml:"var_type,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty"`
}

type RemoteCommandResponse struct {
	HashedId             string                             `json:"remote_command_id" yaml:"remote_command_id"`
	Name                 string                             `json:"name" yaml:"name"`
	Command              string                             `json:"command" yaml:"command"`
	WorkDir              string                             `json:"work_dir" yaml:"work_dir"`
	Timeout              int                                `json:"timeout" yaml:"timeout"`
	EnvironmentVariables []RemoteCommandEnvironmentVariable `json:"environment_variables" yaml:"environment_variables"`
	Status               CommandStatus                      `json:"status" yaml:"status"`
	Creator              RemoteCommandCreator               `json:"creator" yaml:"creator"`
	CreatedAt            time.Time                          `json:"created_at" yaml:"created_at"`
}

type RemoteCommandEnvironmentVariable struct {
	Key   string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
}

type CommandStatus struct {
	Lifecycle        string    `json:"lifecycle" yaml:""`
	CancellationTime time.Time `json:"cancellationTime" yaml:"cancellationTime"`
	ExecutionStart   time.Time `json:"executionStart" yaml:"executionStart"`
	ExecutionEnd     time.Time `json:"executionEnd" yaml:"executionEnd"`
	Error            string    `json:"error" yaml:"error"`
}

type RemoteCommandCreator struct {
	Kind      string `json:"kind"`
	TokenType string `json:"token_type,omitempty"`
	Name      string `json:"name"`
	HashedId  string `json:"hashed_id"`
}

type CacheInvalidation struct {
	Name      string    `json:"name,omitempty" yaml:"name,omitempty"`
	Status    string    `json:"status,omitempty" yaml:"status,omitempty"`
	Objects   []string  `json:"objects,omitempty" yaml:"objects,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty"`
}

type AntivirusScan struct {
	Result        string    `json:"result,omitempty" yaml:"result,omitempty"`
	ScannedFiles  int       `json:"scanned_files,omitempty" yaml:"scanned_files,omitempty"`
	InfectedFiles int       `json:"infected_files,omitempty" yaml:"infected_files,omitempty"`
	Duration      int       `json:"duration,omitempty" yaml:"duration,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty"`
}

type APIUsageTerms struct {
	Terms     string    `json:"terms,omitempty" yaml:"terms,omitempty"`
	Version   string    `json:"version,omitempty" yaml:"version,omitempty"`
	Status    string    `json:"status,omitempty" yaml:"status,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty"`
}

type BackupRequest struct {
	Name       string    `json:"name,omitempty" yaml:"name,omitempty"`
	Kind       string    `json:"kind,omitempty" yaml:"kind,omitempty"`
	Schedule   string    `json:"schedule,omitempty" yaml:"schedule,omitempty"`
	ETA        int       `json:"eta,omitempty" yaml:"eta,omitempty"`
	Components []string  `json:"components,omitempty" yaml:"components,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	KeepUntil  time.Time `json:"keep_until" yaml:"keep_until"`
}

type BackupIterationFlat struct {
	Iteration       string                     `json:"iteration,omitempty" yaml:"iteration,omitempty"`
	ClientName      string                     `json:"client_name,omitempty" yaml:"client_name,omitempty"`
	Reservation     []string                   `json:"reservation,omitempty" yaml:"reservation,omitempty"`
	Components      []BackupIterationComponent `json:"components,omitempty" yaml:"components,omitempty"`
	Status          string                     `json:"status,omitempty" yaml:"status,omitempty"`
	Protection      string                     `json:"protection,omitempty" yaml:"protection,omitempty"`
	EnvironmentName string                     `json:"environment_name,omitempty" yaml:"environment_name,omitempty"`
	ETA             int                        `json:"eta,omitempty" yaml:"eta,omitempty"`
	CreatedAt       time.Time                  `json:"created_at" yaml:"created_at"`
	CompletedAt     time.Time                  `json:"completed_at" yaml:"completed_at"`
}

type BackupIteration struct {
	BackupIterationFlat
	Environment   Environment   `json:"environment" yaml:"environment"`
	BackupRequest BackupRequest `json:"backup_request" yaml:"backup_request"`
}

type Backup struct {
	BackupIteration BackupIterationFlat `json:"backup_iteration" yaml:"backup_iteration"`
	BackupRequest   BackupRequest       `json:"backup_request" yaml:"backup_request"`
}

// BackupIterationResponse ...
type BackupIterationComponent struct {
	Name           string `json:"name,omitempty" yaml:"name,omitempty"`
	BackupSize     int    `json:"backup_size,omitempty" yaml:"backup_size,omitempty"`
	BackupDuration int    `json:"backup_duration,omitempty" yaml:"backup_duration,omitempty"`
	ArchiveKey     string `json:"archive_key,omitempty" yaml:"archive_key,omitempty"`
	IndexKey       string `json:"index_key,omitempty" yaml:"index_key,omitempty"`
	Result         string `json:"result,omitempty" yaml:"result,omitempty"`
}

type RestoreRequestFlat struct {
	Name        string                 `json:"name,omitempty" yaml:"name,omitempty"`
	ETA         int                    `json:"eta,omitempty" yaml:"eta,omitempty"`
	Status      string                 `json:"status,omitempty" yaml:"status,omitempty"`
	Components  []string               `json:"components,omitempty" yaml:"components,omitempty"`
	Results     []RestoreRequestResult `json:"results,omitempty" yaml:"results,omitempty"`
	CreatedAt   time.Time              `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	CompletedAt time.Time              `json:"completed_at,omitempty" yaml:"completed_at,omitempty"`
}

type RestoreRequest struct {
	RestoreRequestFlat
	Initiator       UserMinimal         `json:"initiator" yaml:"initiator"`
	Environment     Environment         `json:"environment" yaml:"environment"`
	BackupIteration BackupIterationFlat `json:"backup_iteration" yaml:"backup_iteration"`
}

// RestoreRequestResponse ...
type RestoreRequestResult struct {
	Name      string    `json:"name,omitempty" yaml:"name,omitempty"`
	Result    string    `json:"result,omitempty" yaml:"result,omitempty"`
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`
}

// RestoreRequestResponse ...
type CWLogResponse struct {
	Timestamp     int64                  `json:"timestamp,omitempty" yaml:"timestamp,omitempty"`
	IngestionTime int64                  `json:"ingestionTime,omitempty" yaml:"ingestionTime,omitempty"`
	LogStreamName string                 `json:"logStreamName,omitempty" yaml:"logStreamName,omitempty"`
	EventId       string                 `json:"eventId,omitempty" yaml:"eventId,omitempty"`
	Log           map[string]interface{} `json:"log,omitempty" yaml:"log,omitempty"`
}

// LogQueryRequest ...
type LogQueryRequest struct {
	Sources   []string  `json:"sources,omitempty" yaml:"sources,omitempty"`
	Filenames []string  `json:"filenames,omitempty" yaml:"filenames,omitempty"`
	Search    string    `json:"search,omitempty" yaml:"search,omitempty"`
	Direction string    `json:"direction,omitempty" yaml:"direction,omitempty"`
	Output    string    `json:"output,omitempty" yaml:"output,omitempty"`
	Limit     int       `json:"limit,omitempty" yaml:"limit,omitempty"`
	Start     time.Time `json:"start,omitempty" yaml:"start,omitempty"`
	End       time.Time `json:"end,omitempty" yaml:"end,omitempty"`
}

type CustomerLogsResponse struct {
	Results []CustomerLogsResult `json:"results" yaml:"results,omitempty"`
}

type CustomerLogsResult struct {
	Stream CustomerLogsStream `json:"stream" yaml:"stream,omitempty"`
	Values [][]string         `json:"values" yaml:"values,omitempty"`
}

type CustomerLogsStream struct {
	Subscription string `json:"subscription" yaml:"subscription,omitempty"`
	Environment  string `json:"environment" yaml:"environment,omitempty"`
	Filename     string `json:"filename,omitempty" yaml:"filename,omitempty"`
	Source       string `json:"source,omitempty" yaml:"source,omitempty"`
	Pod          string `json:"pod,omitempty" yaml:"pod,omitempty"`
}

// CWLogStreamsResponse ...
type CWLogStreamsResponse struct {
	LastEventTimestamp int64  `json:"lastEventTimestamp,omitempty" yaml:"lastEventTimestamp,omitempty"`
	LogStreamName      string `json:"logStreamName,omitempty" yaml:"logStreamName,omitempty"`
}

// LogLabelValuesResponse ...
type LogLabelValuesResponse struct {
	Result []string `json:"result,omitempty" yaml:"result,omitempty"`
}

// UploadResponse
type UploadResponse struct {
	PackageName string `json:"packageName,omitempty" yaml:"packageName,omitempty"`
	BuildID     string `json:"buildId,omitempty" yaml:"buildId,omitempty"`
	BuildName   string `json:"buildName,omitempty" yaml:"buildName,omitempty"`
}

type SyncRequestFlat struct {
	Name            string    `json:"name,omitempty" yaml:"name,omitempty"`
	ETA             int       `json:"eta,omitempty" yaml:"eta,omitempty"`
	Status          string    `json:"status,omitempty" yaml:"status,omitempty"`
	RestoreStrategy string    `json:"results,omitempty" yaml:"results,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	CompletedAt     time.Time `json:"completed_at,omitempty" yaml:"completed_at,omitempty"`
}

type SyncRequest struct {
	RestoreRequestFlat
	Initiator       UserMinimal        `json:"initiator" yaml:"initiator"`
	SrcEnvironment  Environment        `json:"src_environment" yaml:"src_environment"`
	DestEnvironment Environment        `json:"dest_environment" yaml:"dest_environment"`
	RestoreRequest  RestoreRequestFlat `json:"restore_request" yaml:"restore_request"`
	BackupRequest   BackupRequest      `json:"backup_request" yaml:"backup_request"`
}
