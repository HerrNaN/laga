package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumiverse/pulumi-scaleway/sdk/go/scaleway/account"
	"github.com/pulumiverse/pulumi-scaleway/sdk/go/scaleway/iam"
	"github.com/pulumiverse/pulumi-scaleway/sdk/go/scaleway/object"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		region, ok := ctx.GetConfig("scaleway:region")
		if !ok || region == "" {
			region = "fr-par"
		}

		project, err := account.NewProject(ctx, "laga", &account.ProjectArgs{
			Name:        pulumi.String("laga"),
			Description: pulumi.String("laga shopping list app"),
		})
		if err != nil {
			return err
		}

		app, err := iam.NewApplication(ctx, "laga-deploy", &iam.ApplicationArgs{
			Name:        pulumi.String("laga-deploy"),
			Description: pulumi.String("Deployment pipeline for laga"),
		})
		if err != nil {
			return err
		}

		_, err = iam.NewPolicy(ctx, "laga-deploy-policy", &iam.PolicyArgs{
			Name:          pulumi.String("laga-deploy-policy"),
			Description:   pulumi.String("Deploy pipeline: manage containers, serverless SQL, and object storage in laga project"),
			ApplicationId: app.ID(),
			Rules: iam.PolicyRuleArray{
				&iam.PolicyRuleArgs{
					ProjectIds: pulumi.StringArray{
						project.ID(),
					},
					PermissionSetNames: pulumi.StringArray{
						pulumi.String("ContainersFullAccess"),
						pulumi.String("ObjectStorageFullAccess"),
						pulumi.String("ServerlessSQLDatabaseFullAccess"),
					},
					Condition: pulumi.String(""),
				},
			},
		})
		if err != nil {
			return err
		}

		apiKey, err := iam.NewApiKey(ctx, "laga-deploy-key", &iam.ApiKeyArgs{
			ApplicationId:    app.ID(),
			DefaultProjectId: project.ID(),
			Description:      pulumi.String("Deploy pipeline API key for laga"),
			ExpiresAt:        pulumi.String("2027-05-24T00:00:00Z"),
		})
		if err != nil {
			return err
		}

		databaseApp, err := iam.NewApplication(ctx, "laga-database", &iam.ApplicationArgs{
			Name:           pulumi.String("laga-database"),
			Description:    pulumi.String("Laga database connections and schema migrations"),
			OrganizationId: project.OrganizationId,
		})
		if err != nil {
			return err
		}

		databasePolicy, err := iam.NewPolicy(ctx, "laga-database-policy", &iam.PolicyArgs{
			Name:           pulumi.String("laga-database-policy"),
			Description:    pulumi.String("Read/write database access including startup schema migrations"),
			ApplicationId:  databaseApp.ID(),
			OrganizationId: databaseApp.OrganizationId,
			Rules: iam.PolicyRuleArray{
				&iam.PolicyRuleArgs{
					ProjectIds: pulumi.StringArray{project.ID()},
					PermissionSetNames: pulumi.StringArray{
						pulumi.String("ServerlessSQLDatabaseReadWrite"),
					},
					Condition: pulumi.String(""),
				},
			},
		})
		if err != nil {
			return err
		}

		databaseKey, err := iam.NewApiKey(ctx, "laga-database-key", &iam.ApiKeyArgs{
			ApplicationId:    databaseApp.ID(),
			DefaultProjectId: project.ID(),
			Description:      pulumi.String("Laga runtime database API key"),
		}, pulumi.DependsOn([]pulumi.Resource{databasePolicy}))
		if err != nil {
			return err
		}

		infraBucket, err := object.NewBucket(ctx, "laga-pulumi-state", &object.BucketArgs{
			Name:      pulumi.String("laga-pulumi-state"),
			ProjectId: project.ID(),
			Region:    pulumi.String(region),
		})
		if err != nil {
			return err
		}

		ctx.Export("project_id", project.ID())
		ctx.Export("organization_id", app.OrganizationId)
		ctx.Export("access_key", apiKey.AccessKey)
		ctx.Export("secret_key", apiKey.SecretKey)
		ctx.Export("infra_bucket", infraBucket.Name)
		ctx.Export("database_application_id", databaseApp.ID())
		ctx.Export("database_secret_key", pulumi.ToSecret(databaseKey.SecretKey))

		return nil
	})
}
