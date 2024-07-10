package connector

import (
	"context"

	"github.com/SSHcom/privx-sdk-go/api/rolestore"
	"github.com/conductorone/baton-privx/pkg/connector/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

type userBuilder struct {
	client client.PrivXClient
}

func (o *userBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return userResourceType
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (o *userBuilder) List(
	ctx context.Context,
	parentResourceID *v2.ResourceId,
	pToken *pagination.Token,
) (
	[]*v2.Resource,
	string,
	annotations.Annotations,
	error,
) {
	logger := ctxzap.Extract(ctx)
	logger.Debug(
		"Starting call to Users.List",
		zap.String("pToken", pToken.Token),
	)

	offset, err := parsePageToken(pToken.Token)
	if err != nil {
		logger.Error("invalid page token", zap.Error(err))
	}

	privXUsers, nextToken, err := o.client.GetUsers(ctx, offset, pToken.Size)
	if err != nil {
		logger.Debug(
			"Error fetching users",
			zap.Error(err),
		)
		return nil, "", nil, err
	}

	userResources := make([]*v2.Resource, 0)
	for _, user := range privXUsers {
		userCopy := user
		newUserResource, err := userResource(ctx, &userCopy)
		if err != nil {
			return nil, "", nil, err
		}

		userResources = append(userResources, newUserResource)
	}

	return userResources, nextToken, nil, nil
}

// Entitlements always returns an empty slice for users.
func (o *userBuilder) Entitlements(
	_ context.Context,
	resource *v2.Resource,
	_ *pagination.Token,
) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (o *userBuilder) Grants(
	ctx context.Context,
	resource *v2.Resource,
	pToken *pagination.Token,
) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newUserBuilder(client client.PrivXClient) *userBuilder {
	return &userBuilder{client: client}
}

// userResource Converts a PrivX User into a ConductorOne Resource.
func userResource(ctx context.Context, user *rolestore.User) (*v2.Resource, error) {
	createdResource, err := resource.NewUserResource(
		user.FullName,
		userResourceType,
		user.ID,
		[]resource.UserTraitOption{
			resource.WithUserProfile(
				map[string]interface{}{
					"full_name": user.FullName,
					"id":        user.ID,
				},
			),
			resource.WithEmail(user.Email, true),
			resource.WithStatus(v2.UserTrait_Status_STATUS_ENABLED),
		},
	)
	if err != nil {
		return nil, err
	}

	return createdResource, nil
}
