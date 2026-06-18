package connector

import (
	"context"
	"fmt"

	"github.com/SSHcom/privx-sdk-go/api/rolestore"
	"github.com/conductorone/baton-privx/pkg/connector/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

const (
	EntitlementAssigned = "assigned"
)

type roleBuilder struct {
	client client.PrivXClient
}

func (o *roleBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return roleResourceType
}

// List returns all the roles from the upstream service as resource objects.
func (o *roleBuilder) List(
	ctx context.Context,
	parentResourceID *v2.ResourceId,
	opts resource.SyncOpAttrs,
) (
	[]*v2.Resource,
	*resource.SyncOpResults,
	error,
) {
	logger := ctxzap.Extract(ctx)
	pToken := &opts.PageToken

	offset, limit, err := parsePageToken(pToken)
	if err != nil {
		logger.Error("invalid page token", zap.Error(err))
	}

	privXRoles, nextToken, err := o.client.GetRoles(ctx, offset, limit)
	if err != nil {
		logger.Debug("Error fetching roles", zap.Error(err))
		return nil, nil, err
	}

	roleResources := make([]*v2.Resource, 0)
	for _, role := range privXRoles {
		roleCopy := role
		newResource, err := roleResource(ctx, &roleCopy)
		if err != nil {
			return nil, nil, err
		}

		roleResources = append(roleResources, newResource)
	}

	if nextToken == "" {
		return roleResources, nil, nil
	}

	return roleResources, &resource.SyncOpResults{NextPageToken: nextToken}, nil
}

func (o *roleBuilder) Entitlements(
	_ context.Context,
	res *v2.Resource,
	_ resource.SyncOpAttrs,
) ([]*v2.Entitlement, *resource.SyncOpResults, error) {
	entitlements := []*v2.Entitlement{
		entitlement.NewAssignmentEntitlement(
			res,
			EntitlementAssigned,
			entitlement.WithGrantableTo(roleResourceType),
			entitlement.WithDescription(fmt.Sprintf("Has %s role membership", res.DisplayName)),
			entitlement.WithDisplayName(fmt.Sprintf("%s role %s", res.DisplayName, EntitlementAssigned)),
		),
	}
	return entitlements, nil, nil
}

func (o *roleBuilder) Grants(
	ctx context.Context,
	res *v2.Resource,
	opts resource.SyncOpAttrs,
) ([]*v2.Grant, *resource.SyncOpResults, error) {
	logger := ctxzap.Extract(ctx)
	pToken := &opts.PageToken
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
		res.Id.Resource,
		offset,
		limit,
	)
	if err != nil {
		return nil, nil, err
	}

	var roleAssignments []*v2.Grant
	for _, user := range privXUsers {
		roleAssignments = append(
			roleAssignments,
			grant.NewGrant(
				res,
				EntitlementAssigned,
				&v2.ResourceId{
					ResourceType: userResourceType.Id,
					Resource:     user.ID,
				},
			),
		)
	}

	if nextToken == "" {
		return roleAssignments, nil, nil
	}

	return roleAssignments, &resource.SyncOpResults{NextPageToken: nextToken}, nil
}

func (o *roleBuilder) Grant(
	ctx context.Context,
	principal *v2.Resource,
	ent *v2.Entitlement,
) ([]*v2.Grant, annotations.Annotations, error) {
	logger := ctxzap.Extract(ctx)

	if principal.Id.ResourceType != userResourceType.Id {
		logger.Warn(
			"baton-privx: only users can be assigned roles",
			zap.String("principal_type", principal.Id.ResourceType),
			zap.String("principal_id", principal.Id.Resource),
		)
		return nil, nil, fmt.Errorf("baton-privx: only users can be assigned roles")
	}

	err := o.client.GrantRole(
		ctx,
		principal.Id.Resource,
		ent.Resource.Id.Resource,
	)
	if err != nil {
		return nil, nil, err
	}

	newGrant := grant.NewGrant(ent.Resource, EntitlementAssigned, principal.Id)
	return []*v2.Grant{newGrant}, nil, nil
}

func (o *roleBuilder) Revoke(
	ctx context.Context,
	g *v2.Grant,
) (annotations.Annotations, error) {
	logger := ctxzap.Extract(ctx)

	ent := g.Entitlement
	principal := g.Principal

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
		ent.Resource.Id.Resource,
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
