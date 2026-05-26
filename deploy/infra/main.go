package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumiverse/pulumi-scaleway/sdk/go/scaleway/containers"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		region, ok := ctx.GetConfig("scaleway:region")
		if !ok || region == "" {
			region = "fr-par"
		}

		imageTag, ok := ctx.GetConfig("laga:imageTag")
		if !ok || imageTag == "" {
			imageTag = "latest"
		}

		ns, err := containers.NewNamespace(ctx, "laga", &containers.NamespaceArgs{
			Name:      pulumi.String("laga"),
			Region:    pulumi.String(region),
			ProjectId: pulumi.String("0e62dead-55f6-4ddb-8c10-b0a47c5c15b4"),
		})
		if err != nil {
			return err
		}

		container, err := containers.NewContainer(ctx, "laga", &containers.ContainerArgs{
			Name:             pulumi.String("laga"),
			NamespaceId:      ns.ID(),
			Image:            pulumi.Sprintf("ghcr.io/herrnan/laga:%s", imageTag),
			Port:             pulumi.Int(8080),
			Protocol:         pulumi.String("http1"),
			CpuLimit:         pulumi.Int(128),
			MemoryLimitBytes: pulumi.Int(128 * 1024 * 1024),
			MinScale:         pulumi.Int(0),
			MaxScale:         pulumi.Int(1),
			Timeout:          pulumi.Int(5),
			Privacy:          pulumi.String("public"),
			Region:           pulumi.String(region),
		})
		if err != nil {
			return err
		}

		ctx.Export("url", pulumi.Sprintf("https://%s", container.DomainName))
		return nil
	})
}
