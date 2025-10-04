package connector

import (
	"context"

	"github.com/conductorone/baton-minecraft-luckperms/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	resourceSdk "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type userBuilder struct {
	baseURL string
	client  *client.Client
}

func (o *userBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return userResourceType
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (o *userBuilder) List(ctx context.Context, _ *v2.ResourceId, _ *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	users, err := o.client.ListAllUsers(ctx)
	if err != nil {
		return nil, "", nil, err
	}

	resources := []*v2.Resource{}
	for _, u := range users {
		userTraitOptions := []resourceSdk.UserTraitOption{
			resourceSdk.WithStatus(v2.UserTrait_Status_STATUS_ENABLED),
			resourceSdk.WithUserLogin(u.Username),
		}

		userResource, err := resourceSdk.NewUserResource(u.Username, userResourceType, u.UniqueID, userTraitOptions)
		if err != nil {
			return nil, "", nil, err
		}
		resources = append(resources, userResource)
	}

	if len(resources) == 0 {
		return nil, "", nil, nil
	}

	return resources, "", nil, nil
}

// Entitlements always returns an empty slice for users.
func (o *userBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (o *userBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newUserBuilder(baseURL string, client *client.Client) *userBuilder {
	return &userBuilder{
		baseURL: baseURL,
		client:  client,
	}
}
