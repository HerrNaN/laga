package main

import (
	"fmt"
	"net/url"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"github.com/pulumiverse/pulumi-scaleway/sdk/go/scaleway/containers"
	"github.com/pulumiverse/pulumi-scaleway/sdk/go/scaleway/databases"
)

func main() {
	pulumi.Run(deploy)
}

func deploy(ctx *pulumi.Context) error {
	region, ok := ctx.GetConfig("scaleway:region")
	if !ok || region == "" {
		region = "fr-par"
	}

	imageTag, ok := ctx.GetConfig("laga:imageTag")
	if !ok || imageTag == "" {
		imageTag = "latest"
	}

	cfg := config.New(ctx, "laga")
	databaseApplicationID := cfg.Require("databaseApplicationId")
	databaseSecretKey := cfg.RequireSecret("databaseSecretKey")
	const projectID = "0e62dead-55f6-4ddb-8c10-b0a47c5c15b4"

	database, err := databases.NewServerlessDatabase(ctx, "laga", &databases.ServerlessDatabaseArgs{
		Name:      pulumi.String("laga"),
		ProjectId: pulumi.String(projectID),
		Region:    pulumi.String(region),
		MinCpu:    pulumi.Int(0),
		MaxCpu:    pulumi.Int(1),
	}, pulumi.Protect(true))
	if err != nil {
		return err
	}

	databaseURL := pulumi.All(database.Endpoint, databaseSecretKey).ApplyT(func(values []interface{}) (string, error) {
		return databaseConnectionURL(values[0].(string), databaseApplicationID, values[1].(string))
	}).(pulumi.StringOutput)

	ns, err := containers.NewNamespace(ctx, "laga", &containers.NamespaceArgs{
		Name:      pulumi.String("laga"),
		Region:    pulumi.String(region),
		ProjectId: pulumi.String(projectID),
	})
	if err != nil {
		return err
	}

	container, err := containers.NewContainer(ctx, "laga", &containers.ContainerArgs{
		Name:                 pulumi.String("laga"),
		NamespaceId:          ns.ID(),
		Image:                pulumi.Sprintf("ghcr.io/herrnan/laga:%s", imageTag),
		Port:                 pulumi.Int(8080),
		Protocol:             pulumi.String("http1"),
		CpuLimit:             pulumi.Int(128),
		MemoryLimitBytes:     pulumi.Int(128 * 1024 * 1024),
		HttpsConnectionsOnly: pulumi.Bool(true),
		MinScale:             pulumi.Int(0),
		MaxScale:             pulumi.Int(1),
		Timeout:              pulumi.Int(5),
		Privacy:              pulumi.String("public"),
		Region:               pulumi.String(region),
		SecretEnvironmentVariables: pulumi.StringMap{
			"DATABASE_URL": databaseURL,
		},
	})
	if err != nil {
		return err
	}

	ctx.Export("url", pulumi.Sprintf("https://%s", container.PublicEndpoint))
	ctx.Export("database_endpoint", database.Endpoint)
	return nil
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
