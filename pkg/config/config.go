package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	BaseUrlField = field.StringField(
		"base-url",
		field.WithRequired(true),
		field.WithDisplayName("Base URL"),
		field.WithPlaceholder("https://your-privx-instance.example.com"),
		field.WithDescription("The hostname (URL) for your PrivX instance"),
	)
	ApiClientIdField = field.StringField(
		"api-client-id",
		field.WithRequired(true),
		field.WithDisplayName("API Client ID"),
		field.WithPlaceholder("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"),
		field.WithDescription("The API Client ID (a UUID.)"),
	)
	ApiClientSecretField = field.StringField(
		"api-client-secret",
		field.WithRequired(true),
		field.WithIsSecret(true),
		field.WithDisplayName("API Client Secret"),
		field.WithDescription("The API Client Secret (a base64 string.)"),
	)
	OauthClientIdField = field.StringField(
		"oauth-client-id",
		field.WithRequired(true),
		field.WithDisplayName("OAuth Client ID"),
		field.WithPlaceholder("privx-external"),
		field.WithDescription("The OAuth Client ID (e.g. \"privx-external\".)"),
	)
	OauthClientSecretField = field.StringField(
		"oauth-client-secret",
		field.WithRequired(true),
		field.WithIsSecret(true),
		field.WithDisplayName("OAuth Client Secret"),
		field.WithDescription("The OAuth Client Secret (a base64 string.)"),
	)

	ConfigurationFields = []field.SchemaField{
		BaseUrlField,
		ApiClientIdField,
		ApiClientSecretField,
		OauthClientIdField,
		OauthClientSecretField,
	}
)

//go:generate go run ./gen
var Config = field.NewConfiguration(
	ConfigurationFields,
	field.WithConnectorDisplayName("PrivX"),
	field.WithIconUrl("/static/app-icons/privx.svg"),
	field.WithHelpUrl("/docs/baton/privx"),
)
