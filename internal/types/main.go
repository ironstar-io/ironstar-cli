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
	Environment  Environment   `yaml:"environment,omitempty"`
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
	HashedID       string `json:"environment_id,omitempty" yaml:"environment_id,omitempty"`
	Class          string `json:"class,omitempty" yaml:"class,omitempty"`
	DNSName        string `json:"dns_name,omitempty" yaml:"dns_name,omitempty"`
	Name           string `json:"name,omitempty" yaml:"name,omitempty"`
	Trigger        string `json:"trigger,omitempty" yaml:"trigger,omitempty"`
	UpdateStrategy string `json:"update_strategy,omitempty" yaml:"update_strategy,omitempty"`
}

// UploadResponse
type UploadResponse struct {
	PackageName string `json:"packageName,omitempty" yaml:"packageName,omitempty"`
	BuildID     string `json:"buildId,omitempty" yaml:"buildId,omitempty"`
	BuildName   string `json:"buildName,omitempty" yaml:"buildName,omitempty"`
}
