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
	Project      Project      `yaml:"project,omitempty"`
	Subscription Subscription `yaml:"subscription,omitempty"`
}

type Project struct {
	Name string `yaml:"name,omitempty"`
	Path string `yaml:"path,omitempty"`
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
