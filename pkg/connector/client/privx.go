package client

import (
	"github.com/SSHcom/privx-sdk-go/restapi"
	"github.com/SSHcom/privx-sdk-go/oauth"
)

// 1. Create Authorizer and Access Token Provider
func authorize() restapi.Authorizer {
	auth := restapi.New(
		/* use restapi options to config http */
		/* the options can be referred from SDK Configuration providers section below*/
		restapi.UseConfigFile("config.toml"),
		restapi.UseEnvironment(),
	)

	return oauth.With(
		auth,
		// 1. Use config file option to configure authorizer
		oauth.UseConfigFile("config.toml"),
		// 2. Use environment variables option to configure authorizer
		oauth.UseEnvironment(),
		// 3. Use oauth options to configure authorizer
		oauth.Access(/* ... */),
		oauth.Secret(/* ... */),
	)
}

// 2. Create HTTP transport for PrivX API
func curl() restapi.Connector {
	return restapi.New(
		restapi.Auth(authorize()),
	)
}

// 3. Create rolestore instance with API client/connector
roleStore := rolestore.New(curl())