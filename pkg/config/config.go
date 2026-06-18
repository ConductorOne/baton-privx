package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	BaseUrlField = field.StringField(
		"base-url",
		field.WithDisplayName("Base URL"),
		field.WithPlaceholder("Your PrivX base URL"),
	)
	ClientIdField = field.StringField(
		"client-id",
		field.WithDisplayName("Client ID"),
		field.WithPlaceholder("Your PrivX client ID"),
	)
	ClientSecretField = field.StringField(
		"client-secret",
		field.WithIsSecret(true),
		field.WithDisplayName("Client secret"),
		field.WithPlaceholder("Your PrivX client secret"),
	)
	OauthClientIdField = field.StringField(
		"oauth-client-id",
		field.WithDisplayName("OAuth client ID"),
		field.WithPlaceholder("Your PrivX OAuth client ID"),
	)
	OauthClientSecretField = field.StringField(
		"oauth-client-secret",
		field.WithIsSecret(true),
		field.WithDisplayName("OAuth client secret"),
		field.WithPlaceholder("Your PrivX OAuth client secret"),
	)

	ConfigurationFields = []field.SchemaField{
		BaseUrlField,
		ClientIdField,
		ClientSecretField,
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
	field.WithFieldGroups([]field.SchemaFieldGroup{
		{
			Name:        "privx-group-oauth",
			DisplayName: "OAuth",
			Fields:      []field.SchemaField{BaseUrlField, OauthClientIdField, OauthClientSecretField},
		},
		{
			Name:        "privx-group-client-secret",
			DisplayName: "Client secret",
			Fields:      []field.SchemaField{BaseUrlField, ClientIdField, ClientSecretField},
		},
	}),
)
