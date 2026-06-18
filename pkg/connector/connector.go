package connector

import (
	"context"
	"fmt"
	"io"

	cfg "github.com/conductorone/baton-privx/pkg/config"
	"github.com/conductorone/baton-privx/pkg/connector/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/cli"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
)

type Config struct {
	BaseUrl           string
	PrivxClientId     string
	PrivxClientSecret string
	OAuthClientID     string
	OAuthClientSecret string
}

type Connector struct {
	client client.PrivXClient
}

// ResourceSyncers returns a ResourceSyncerV2 for each resource type that should be synced from the upstream service.
func (d *Connector) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncerV2 {
	return []connectorbuilder.ResourceSyncerV2{
		newUserBuilder(d.client),
		newRoleBuilder(d.client),
	}
}

// Asset takes an input AssetRef and attempts to fetch it using the connector's authenticated http client
// It streams a response, always starting with a metadata object, following by chunked payloads for the asset.
func (d *Connector) Asset(ctx context.Context, asset *v2.AssetRef) (string, io.ReadCloser, error) {
	return "", nil, nil
}

// Metadata returns metadata about the connector.
func (d *Connector) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "PrivX",
		Description: "Baton connector for PrivX",
	}, nil
}

// Validate is called to ensure that the connector is properly configured. It should exercise any API credentials
// to be sure that they are valid.
func (d *Connector) Validate(ctx context.Context) (annotations.Annotations, error) {
	err := d.client.Verify(ctx)
	if err != nil {
		return nil, fmt.Errorf("privx-connector: failed to validate client credentials: %w", err)
	}

	return nil, nil
}

// NewWithOAuth returns a new instance of the connector using OAuth authentication.
func NewWithOAuth(ctx context.Context, baseUrl, oAuthClientID, oAuthClientSecret string) (*Connector, error) {
	privXClient, err := client.NewPrivXClientWithOAuth(ctx, baseUrl, oAuthClientID, oAuthClientSecret)
	if err != nil {
		return nil, err
	}
	return &Connector{client: *privXClient}, nil
}

// NewWithClientSecret returns a new instance of the connector using client secret authentication.
func NewWithClientSecret(ctx context.Context, baseUrl, clientId, clientSecret string) (*Connector, error) {
	privXClient, err := client.NewPrivXClientWithClientSecret(ctx, baseUrl, clientId, clientSecret)
	if err != nil {
		return nil, err
	}
	return &Connector{client: *privXClient}, nil
}

// NewLambdaConnector returns a new ConnectorBuilderV2 for use with RunConnector / Lambda.
func NewLambdaConnector(ctx context.Context, ac *cfg.Privx, opts *cli.ConnectorOpts) (connectorbuilder.ConnectorBuilderV2, []connectorbuilder.Opt, error) {
	var c *Connector
	var err error

	switch opts.SelectedAuthMethod {
	case "privx-group-oauth":
		c, err = NewWithOAuth(ctx, ac.BaseUrl, ac.OauthClientId, ac.OauthClientSecret)
	case "privx-group-client-secret":
		c, err = NewWithClientSecret(ctx, ac.BaseUrl, ac.PrivxClientId, ac.PrivxClientSecret)
	default:
		return nil, nil, fmt.Errorf("baton-privx: unknown auth method: %s", opts.SelectedAuthMethod)
	}
	if err != nil {
		return nil, nil, err
	}
	return c, nil, nil
}
