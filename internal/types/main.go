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

// Project is a singular entry of a project name and path used in global config
type ProjectConfig struct {
	Name         string       `yaml:"name,omitempty"`
	Path         string       `yaml:"path,omitempty"`
	Login        string       `yaml:"login,omitempty"`
	Subscription Subscription `yaml:"subscription,omitempty"`
}

// Subscription
type Subscription struct {
	HashedID        string `json:"subscription_id,omitempty" example:"98hreHs"`
	ApplicationType string `json:"application_type,omitempty" example:"drupal"`
	Ref             string `json:"ref,omitempty" example:"au1999"`
	Alias           string `json:"alias,omitempty" example:"umami-food-blog"`
}
