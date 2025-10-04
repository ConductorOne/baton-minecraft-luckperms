package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	// Add the SchemaFields for the Config.
	configField         = field.StringField("configField")
	ConfigurationFields = []field.SchemaField{configField, Address, Port}

	Address = field.StringField(
		"address",
		field.WithDisplayName("Address"),
		field.WithIsSecret(false),
		field.WithDescription("The URL of the LuckPerms REST endpoint."),
		field.WithRequired(true),
	)
	Port = field.StringField(
		"port",
		field.WithDisplayName("Port"),
		field.WithIsSecret(false),
		field.WithDescription("The Port of the LuckPerms REST endpoint."),
		field.WithRequired(true),
	)

	// FieldRelationships defines relationships between the ConfigurationFields that can be automatically validated.
	// For example, a username and password can be required together, or an access token can be
	// marked as mutually exclusive from the username password pair.
	FieldRelationships = []field.SchemaFieldRelationship{
		{
			Kind: field.RequiredTogether,
			Fields: []field.SchemaField{
				Address,
				Port,
			},
		},
	}
)

//go:generate go run -tags=generate ./gen
var Config = field.NewConfiguration(
	ConfigurationFields,
	field.WithConstraints(FieldRelationships...),
	field.WithConnectorDisplayName("Minecraft Luckperms"),
)
