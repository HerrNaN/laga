package main

import (
	"net/url"
	"sync"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func TestDatabaseConnectionURL(t *testing.T) {
	for _, endpoint := range []string{
		"postgres://database.example.com:5432/laga",
		"postgresql://database.example.com:5432/laga?connect_timeout=10&sslmode=disable",
	} {
		t.Run(endpoint, func(t *testing.T) {
			const applicationID = "database-application"
			const secretKey = "secret:@/?#%+ key"
			connectionURL, err := databaseConnectionURL(endpoint, applicationID, secretKey)
			if err != nil {
				t.Fatal(err)
			}
			parsed, err := url.Parse(connectionURL)
			if err != nil {
				t.Fatal(err)
			}
			password, ok := parsed.User.Password()
			if parsed.User.Username() != applicationID || !ok || password != secretKey {
				t.Fatal("database credentials did not round-trip")
			}
			if parsed.Host != "database.example.com:5432" || parsed.Path != "/laga" {
				t.Fatal("database address changed")
			}
			if parsed.Query().Get("sslmode") != "require" {
				t.Fatal("TLS is not required")
			}
			original, _ := url.Parse(endpoint)
			if parsed.Query().Get("connect_timeout") != original.Query().Get("connect_timeout") {
				t.Fatal("endpoint query parameters were lost")
			}
		})
	}
}

func TestDatabaseConnectionURLRejectsInvalidEndpoint(t *testing.T) {
	for _, endpoint := range []string{"", "database.example.com", "https://database.example.com/laga", "postgres://", "postgres://%invalid"} {
		t.Run(endpoint, func(t *testing.T) {
			if _, err := databaseConnectionURL(endpoint, "application", "secret"); err == nil {
				t.Fatal("expected an invalid endpoint error")
			}
		})
	}
}

type deploymentMocks struct {
	resources sync.Map
}

func (m *deploymentMocks) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	return args.Args, nil
}

func (m *deploymentMocks) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	m.resources.Store(args.TypeToken, args)
	outputs := args.Inputs.Copy()
	switch args.TypeToken {
	case "scaleway:databases/serverlessDatabase:ServerlessDatabase":
		outputs["endpoint"] = resource.NewStringProperty("postgres://database.example.com:5432/laga")
	case "scaleway:containers/container:Container":
		outputs["publicEndpoint"] = resource.NewStringProperty("laga.example.com")
	}
	return args.Name + "-id", outputs, nil
}

func TestDeployWiresSecretDatabaseURL(t *testing.T) {
	t.Setenv("PULUMI_CONFIG", `{"laga:databaseApplicationId":"database-application","laga:databaseSecretKey":"database-secret"}`)
	mocks := &deploymentMocks{}
	if err := pulumi.RunErr(deploy, pulumi.WithMocks("laga-infra", "test", mocks)); err != nil {
		t.Fatal(err)
	}

	dbValue, ok := mocks.resources.Load("scaleway:databases/serverlessDatabase:ServerlessDatabase")
	if !ok {
		t.Fatal("serverless database was not registered")
	}
	db := dbValue.(pulumi.MockResourceArgs)
	if db.Inputs["minCpu"].NumberValue() != 0 || db.Inputs["maxCpu"].NumberValue() != 1 {
		t.Fatal("expected database scaling from zero to one CPU")
	}
	if db.RegisterRPC == nil || !db.RegisterRPC.GetProtect() {
		t.Fatal("database is not protected against deletion")
	}

	containerValue, ok := mocks.resources.Load("scaleway:containers/container:Container")
	if !ok {
		t.Fatal("container was not registered")
	}
	container := containerValue.(pulumi.MockResourceArgs)
	secrets := container.Inputs["secretEnvironmentVariables"]
	if !secrets.IsSecret() {
		t.Fatal("database environment variables are not marked secret")
	}
	databaseURL := secrets.SecretValue().Element.ObjectValue()["DATABASE_URL"].StringValue()
	expected, err := databaseConnectionURL("postgres://database.example.com:5432/laga", "database-application", "database-secret")
	if err != nil {
		t.Fatal(err)
	}
	if databaseURL != expected {
		t.Fatal("container database URL does not use the provisioned endpoint and runtime credentials")
	}
	if container.Inputs["region"].StringValue() != db.Inputs["region"].StringValue() {
		t.Fatal("container and database regions differ")
	}
}
