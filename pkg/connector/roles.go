package connector

import (
	"context"
	"fmt"

	"github.com/SSHcom/privx-sdk-go/api/rolestore"
	"github.com/conductorone/baton-privx/pkg/connector/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

const (
	// TODO MARCOS 1.0 are these values going to be used?
	RoleMembershipPermanent  = "permanent"
	RoleMembershipRestricted = "restricted"
	RoleMembershipFloating   = "floating"
	EntitlementAssigned      = "assigned"
)

type roleBuilder struct {
	client client.PrivXClient
}

func (o *roleBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return roleResourceType
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (o *roleBuilder) List(
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

	offset, limit, err := parsePageToken(pToken)
	if err != nil {
		logger.Error("invalid page token", zap.Error(err))
	}

	privXRoles, nextToken, err := o.client.GetRoles(ctx, offset, limit)
	if err != nil {
		logger.Debug("Error fetching users", zap.Error(err))
		return nil, "", nil, err
	}

	roleResources := make([]*v2.Resource, 0)
	for _, role := range privXRoles {
		roleCopy := role
		newResource, err := roleResource(ctx, &roleCopy)
		if err != nil {
			return nil, "", nil, err
		}

		roleResources = append(roleResources, newResource)
	}

	return roleResources, nextToken, nil, nil
}

// Entitlements always returns an empty slice for users.
func (o *roleBuilder) Entitlements(
	_ context.Context,
	resource *v2.Resource,
	_ *pagination.Token,
) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	entitlements := []*v2.Entitlement{
		entitlement.NewAssignmentEntitlement(
			resource,
			EntitlementAssigned,
			entitlement.WithGrantableTo(roleResourceType),
			entitlement.WithDescription(fmt.Sprintf("Has %s role membership", resource.DisplayName)),
			entitlement.WithDisplayName(fmt.Sprintf("%s role %s", resource.DisplayName, EntitlementAssigned)),
		),
	}
	return entitlements, "", nil, nil
}

func (o *roleBuilder) Grants(
	ctx context.Context,
	resource *v2.Resource,
	pToken *pagination.Token,
) ([]*v2.Grant, string, annotations.Annotations, error) {
	logger := ctxzap.Extract(ctx)
	logger.Debug(
		"Starting call to Roles.Grants",
		zap.String("pToken", pToken.Token),
	)

	offset, limit, err := parsePageToken(pToken)
	if err != nil {
		logger.Error("invalid page token", zap.Error(err))
	}

	privXUsers, nextToken, err := o.client.GetUsersForRole(
		ctx,
		resource.Id.Resource,
		offset,
		limit,
	)
	if err != nil {
		return nil, "", nil, err
	}

	var roleAssignments []*v2.Grant
	for _, user := range privXUsers {
		roleAssignments = append(
			roleAssignments,
			grant.NewGrant(
				resource,
				EntitlementAssigned,
				&v2.ResourceId{
					ResourceType: userResourceType.Id,
					Resource:     user.ID,
				},
			),
		)
	}

	return roleAssignments, nextToken, nil, nil
}

func (o *roleBuilder) Grant(
	ctx context.Context,
	principal *v2.Resource,
	entitlement *v2.Entitlement,
) (annotations.Annotations, error) {
	logger := ctxzap.Extract(ctx)

	if principal.Id.ResourceType != userResourceType.Id {
		logger.Warn(
			"baton-privx: only users can be assigned roles",
			zap.String("principal_type", principal.Id.ResourceType),
			zap.String("principal_id", principal.Id.Resource),
		)
		return nil, fmt.Errorf("baton-privx: only users can be assigned roles")
	}

	err := o.client.GrantRole(
		ctx,
		principal.Id.Resource,
		entitlement.Resource.Id.Resource,
	)
	return nil, err
}

func (o *roleBuilder) Revoke(
	ctx context.Context,
	grant *v2.Grant,
) (annotations.Annotations, error) {
	logger := ctxzap.Extract(ctx)

	entitlement := grant.Entitlement
	principal := grant.Principal

	if principal.Id.ResourceType != userResourceType.Id {
		logger.Warn(
			"baton-privx: only users can have role assignment revoked",
			zap.String("principal_type", principal.Id.ResourceType),
			zap.String("principal_id", principal.Id.Resource),
		)
		return nil, fmt.Errorf("baton-privx: only users can have role assignment revoked")
	}

	err := o.client.RevokeRole(
		ctx,
		principal.Id.Resource,
		entitlement.Resource.Id.Resource,
	)
	return nil, err
}

func newRoleBuilder(client client.PrivXClient) *roleBuilder {
	return &roleBuilder{client: client}
}

func roleResource(ctx context.Context, role *rolestore.Role) (*v2.Resource, error) {
	createdResource, err := resource.NewRoleResource(
		role.Name,
		roleResourceType,
		role.ID,
		[]resource.RoleTraitOption{
			resource.WithRoleProfile(map[string]interface{}{
				"name": role.Name,
			}),
		},
	)
	if err != nil {
		return nil, err
	}

	return createdResource, nil
}
