package main

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-sdk/pkg/field"
	"github.com/spf13/viper"
)

var (
	baseUrlField = field.StringField(
		"base-url",
		field.WithDescription("The hostname (URL) for your PrivX instance. ($BATON_BASE_URL)"),
	)
	apiClientIdField = field.StringField(
		"api-client-id",
		field.WithDescription("The API Client ID (a UUID.) ($BATON_API_CLIENT_ID)"),
	)
	apiClientSecretField = field.StringField(
		"api-client-secret",
		field.WithDescription("The API Client Secret (a base64 string.) ($BATON_API_CLIENT_SECRET)"),
	)
	oauthClientIdField = field.StringField(
		"oauth-client-id",
		field.WithDescription("The OAuth Client ID (e.g. \"privx-external\".) ($BATON_OAUTH_CLIENT_ID)"),
	)
	oauthClientSecretField = field.StringField(
		"oauth-client-secret",
		field.WithDescription("The OAuth Client Secret (a base64 string.) ($BATON_OAUTH_CLIENT_SECRET)"),
	)
)

// configurationFields defines the external configuration required for the connector to run.
var configurationFields = []field.SchemaField{
	apiClientIdField,
	apiClientSecretField,
	baseUrlField,
	oauthClientIdField,
	oauthClientSecretField,
}

// validateConfig is run after the configuration is loaded, and should return an error if it isn't valid.
func validateConfig(ctx context.Context, v *viper.Viper) error {
	if v.GetString(baseUrlField.FieldName) == "" {
		return fmt.Errorf("base-url is required")
	}
	if v.GetString(apiClientIdField.FieldName) == "" {
		return fmt.Errorf("api-client-id is required")
	}
	if v.GetString(apiClientSecretField.FieldName) == "" {
		return fmt.Errorf("api-client-secret is required")
	}
	if v.GetString(oauthClientIdField.FieldName) == "" {
		return fmt.Errorf("oauth-client-id is required")
	}
	if v.GetString(oauthClientSecretField.FieldName) == "" {
		return fmt.Errorf("oauth-client-secret is required")
	}
	return nil
}
