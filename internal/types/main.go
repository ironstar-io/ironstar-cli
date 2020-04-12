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
	RedirectEndpoint string    `json:"redirect_endpoint"` // If this is set, user is MFA registered
	Expiry           time.Time `json:"expiry"`
}

type ProjectConfig struct {
	Version      string       `yaml:"version,omitempty"`
	Subscription Subscription `yaml:"subscription,omitempty"`
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

// BuildsResponse
type BuildsResponse struct {
	HashedID  string    `json:"build_id,omitempty" yaml:"build_id,omitempty"`
	Status    string    `json:"status,omitempty" yaml:"status,omitempty"`
	CreatedBy string    `json:"created_by,omitempty" yaml:"created_by,omitempty"`
	RunningIn string    `json:"running_in,omitempty" yaml:"running_in,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty"`
}
