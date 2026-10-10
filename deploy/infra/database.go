package main

import (
	"fmt"
	"net/url"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"github.com/pulumiverse/pulumi-scaleway/sdk/go/scaleway/databases"
)

func createDatabase(ctx *pulumi.Context, region string, projectID pulumi.StringInput) (pulumi.StringOutput, error) {
	cfg := config.New(ctx, "laga")
	databaseApplicationID := cfg.Require("databaseApplicationId")
	databaseSecretKey := cfg.RequireSecret("databaseSecretKey")

	database, err := databases.NewServerlessDatabase(ctx, "laga", &databases.ServerlessDatabaseArgs{
		Name:      pulumi.String("laga"),
		ProjectId: projectID,
		Region:    pulumi.String(region),
		MinCpu:    pulumi.Int(0),
		MaxCpu:    pulumi.Int(1),
	}, pulumi.Protect(true))
	if err != nil {
		return pulumi.StringOutput{}, err
	}

	databaseURL := pulumi.All(database.Endpoint, databaseSecretKey).ApplyT(func(values []interface{}) (string, error) {
		return databaseConnectionURL(values[0].(string), databaseApplicationID, values[1].(string))
	}).(pulumi.StringOutput)

	ctx.Export("database_endpoint", database.Endpoint)
	return databaseURL, nil
}

func databaseConnectionURL(endpoint, applicationID, secretKey string) (string, error) {
	connectionURL, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("parse database endpoint: %w", err)
	}
	if (connectionURL.Scheme != "postgres" && connectionURL.Scheme != "postgresql") || connectionURL.Host == "" {
		return "", fmt.Errorf("invalid PostgreSQL database endpoint")
	}
	connectionURL.User = url.UserPassword(applicationID, secretKey)
	query := connectionURL.Query()
	query.Set("sslmode", "require")
	connectionURL.RawQuery = query.Encode()
	return connectionURL.String(), nil
}
