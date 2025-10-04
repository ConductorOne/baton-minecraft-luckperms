package connector

import (
	"context"

	"github.com/conductorone/baton-minecraft-luckperms/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	resourceSdk "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type groupBuilder struct {
	baseURL string
	client  *client.Client
}

func (o *groupBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return groupResourceType
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (o *groupBuilder) List(ctx context.Context, _ *v2.ResourceId, _ *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	groups, err := o.client.ListAllGroups(ctx)
	if err != nil {
		return nil, "", nil, err
	}
	resources := make([]*v2.Resource, 0, len(groups))
	for _, group := range groups {
		data := make(map[string]any)
		for _, node := range group.Nodes {
			data[node.Key] = map[string]any{
				"value": node.Value,
				"type":  node.Type,
				//"context": node.Context,
			}
		}
		data["metadata"] = group.Metadata

		groupTraitOptions := []resourceSdk.GroupTraitOption{
			resourceSdk.WithGroupProfile(data),
		}
		groupResource, err := resourceSdk.NewGroupResource(group.DisplayName, groupResourceType, group.Name, groupTraitOptions)
		if err != nil {
			return nil, "", nil, err
		}
		resources = append(resources, groupResource)
	}

	if len(resources) == 0 {
		return nil, "", nil, nil
	}

	return resources, "", nil, nil
}

// Entitlements always returns an empty slice for users.
func (o *groupBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	opts := []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDisplayName(resource.Id.Resource + " Member"),
	}
	entitlementItem := entitlement.NewAssignmentEntitlement(resource, "member", opts...)

	return []*v2.Entitlement{entitlementItem}, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (o *groupBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newGroupBuilder(baseURL string, client *client.Client) *groupBuilder {
	return &groupBuilder{
		baseURL: baseURL,
		client:  client,
	}
}
