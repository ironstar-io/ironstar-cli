package types

import (
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
	CreatedAt time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty"`
}

// Build
type Build struct {
	HashedID   string       `json:"build_id,omitempty" yaml:"build_id,omitempty"`
	Name       string       `json:"name,omitempty" yaml:"name,omitempty"`
	Ref        string       `json:"ref,omitempty" yaml:"ref,omitempty"`
	Status     string       `json:"status,omitempty" yaml:"status,omitempty"`
	CreatedBy  string       `json:"created_by,omitempty" yaml:"created_by,omitempty"`
	RunningIn  string       `json:"running_in,omitempty" yaml:"running_in,omitempty"`
	CreatedAt  time.Time    `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	Deployment []Deployment `json:"deployments,omitempty" yaml:"deployments,omitempty"`
}

// Deployment
type Deployment struct {
	HashedID       string      `json:"deployment_id,omitempty" yaml:"deployment_id,omitempty"`
	Name           string      `json:"name,omitempty" yaml:"name,omitempty"`
	AppStatus      string      `json:"app_status,omitempty" yaml:"app_status,omitempty"`
	AdminSvcStatus string      `json:"admin_svc_status,omitempty" yaml:"admin_svc_status,omitempty"`
	BuildID        string      `json:"build_token,omitempty" yaml:"build_token,omitempty"`
	CreatedAt      time.Time   `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	Environment    Environment `json:"environment,omitempty" yaml:"environment,omitempty"`
	Build          BuildFlat   `json:"build,omitempty" yaml:"build,omitempty"`
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
}

type BackupRequest struct {
	Name       string    `json:"name,omitempty" yaml:"name,omitempty"`
	Kind       string    `json:"kind,omitempty" yaml:"kind,omitempty"`
	Schedule   string    `json:"schedule,omitempty" yaml:"schedule,omitempty"`
	ETA        int       `json:"eta,omitempty" yaml:"eta,omitempty"`
	Components []string  `json:"components,omitempty" yaml:"components,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty"`
}

type BackupIterationFlat struct {
	Iteration   string                     `json:"iteration,omitempty" yaml:"iteration,omitempty"`
	ClientName  string                     `json:"client_name,omitempty" yaml:"client_name,omitempty"`
	Reservation []string                   `json:"reservation,omitempty" yaml:"reservation,omitempty"`
	Components  []BackupIterationComponent `json:"components,omitempty" yaml:"components,omitempty"`
	Status      string                     `json:"status,omitempty" yaml:"status,omitempty"`
	Protection  string                     `json:"protection,omitempty" yaml:"protection,omitempty"`
	ETA         int                        `json:"eta,omitempty" yaml:"eta,omitempty"`
	CreatedAt   time.Time                  `json:"created_at" yaml:"created_at"`
	CompletedAt time.Time                  `json:"completed_at" yaml:"completed_at"`
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

// UploadResponse
type UploadResponse struct {
	PackageName string `json:"packageName,omitempty" yaml:"packageName,omitempty"`
	BuildID     string `json:"buildId,omitempty" yaml:"buildId,omitempty"`
	BuildName   string `json:"buildName,omitempty" yaml:"buildName,omitempty"`
}
