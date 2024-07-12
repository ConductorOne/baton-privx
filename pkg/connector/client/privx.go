package client

import (
	"context"
	"strconv"
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/rolestore"
	"github.com/SSHcom/privx-sdk-go/oauth"
	"github.com/SSHcom/privx-sdk-go/restapi"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

type PrivXClient struct {
	Authorizer restapi.Authorizer
	RoleStore  rolestore.RoleStore
}

func NewPrivXClient(
	ctx context.Context,
	baseUrl string,
	apiClientId string,
	apiClientSecret string,
	oauthClientId string,
	oauthClientSecret string,
) (*PrivXClient, error) {
	baseUrl = strings.Trim(baseUrl, "/")

	authorizer := oauth.With(
		restapi.New(
			restapi.BaseURL(baseUrl),
		),
		oauth.Access(apiClientId),
		oauth.Secret(apiClientSecret),
		oauth.Digest(oauthClientId, oauthClientSecret),
	)

	roleStore := rolestore.New(
		restapi.New(
			restapi.Auth(authorizer),
			restapi.BaseURL(baseUrl),
		),
	)

	return &PrivXClient{
		Authorizer: authorizer,
		RoleStore:  *roleStore,
	}, nil
}

func getNextToken(start, found, pageSize int) string {
	if found < pageSize {
		return ""
	}

	return strconv.FormatInt(int64(start+found), 10)
}

// Verify fetches an API access token to verify that the client credentials are valid.
func (c *PrivXClient) Verify(ctx context.Context) error {
	logger := ctxzap.Extract(ctx)
	logger.Debug("calling AccessToken with client credentials")

	_, err := c.Authorizer.AccessToken()
	if err != nil {
		logger.Error(
			"could not fetch an API access token with client credentials",
			zap.Error(err),
		)
		return err
	}
	return nil
}

// GetUsers uses pagination to get a list of users from the global list. Returns
// ([]user, string, error) tuple that represents the fetched list of users, the
// next pagination token, and potentially any errors. The next pagination token
// is a string so that we can use `""` as a signal that there are no more pages.
func (c *PrivXClient) GetUsers(
	ctx context.Context,
	offset int,
	limit int,
) (
	[]rolestore.User,
	string,
	error,
) {
	privXUsers, err := c.RoleStore.SearchUsers(
		offset,
		limit,
		"",
		"",
		rolestore.UserSearchObject{},
	)
	if err != nil {
		return nil, "", err
	}

	nextToken := getNextToken(offset, len(privXUsers), limit)

	return privXUsers, nextToken, nil
}

func (c *PrivXClient) GetRoles(
	ctx context.Context,
	offset int,
	limit int,
) (
	[]rolestore.Role,
	string,
	error,
) {
	privXRoles, err := c.RoleStore.Roles(
		offset,
		limit,
		"",
		"",
	)
	if err != nil {
		return nil, "", err
	}

	nextToken := getNextToken(offset, len(privXRoles), limit)

	return privXRoles, nextToken, nil
}

func (c *PrivXClient) GetUsersForRole(
	ctx context.Context,
	roleId string,
	offset int,
	limit int,
) (
	[]rolestore.User,
	string,
	error,
) {
	privXRoles, err := c.RoleStore.GetRoleMembers(
		roleId,
		offset,
		limit,
		"",
		"",
	)
	if err != nil {
		return nil, "", err
	}

	nextToken := getNextToken(offset, len(privXRoles), limit)

	return privXRoles, nextToken, nil
}

// GrantRole fetches the list of roles for a given user and appends the
// specified role to that list. NOTE: the fetch and put are _not_ atomic and
// can cause race conditions.
func (c *PrivXClient) GrantRole(ctx context.Context, userId, roleId string) error {
	return c.RoleStore.GrantUserRole(userId, roleId)
}

// RevokeRole fetches the list of roles for a given user and removes the
// specified role from that list. NOTE: the fetch and put are _not_ atomic and
// can cause race conditions.
func (c *PrivXClient) RevokeRole(ctx context.Context, userId, roleId string) error {
	return c.RoleStore.RevokeUserRole(userId, roleId)
}
