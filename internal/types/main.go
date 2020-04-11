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

type AuthLoginBody struct {
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
	HashedID        string `yaml:"subscription_id,omitempty" example:"98hreHs"`
	Alias           string `yaml:"alias,omitempty" example:"umami-food-blog"`
	Ref             string `yaml:"ref,omitempty" example:"au1999"`
	ApplicationType string `yaml:"application_type,omitempty" example:"drupal"`
}
