package main

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/appengine"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"github.com/pulumiverse/pulumi-sentry/sdk/go/sentry"
	"github.com/upstash/pulumi-upstash/sdk/go/upstash"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		conf := config.New(ctx, "")

		if err := setupGCPProject(ctx, conf); err != nil {
			return err
		}

		if err := setupRedis(ctx, conf); err != nil {
			return err
		}

		project := conf.Get("sentry-project")
		org := conf.Get("sentry-org")
		team := conf.Get("sentry-team")
		if project != "" && org != "" && team != "" {
			if err := setupSentryProject(ctx, project, org, team); err != nil {
				return err
			}
		}

		return nil
	})
}

func setupRedis(ctx *pulumi.Context, conf *config.Config) error {
	redis, err := upstash.NewRedisDatabase(ctx, conf.Get("redisdb"), &upstash.RedisDatabaseArgs{
		DatabaseName: pulumi.String(conf.Get("redisdb")),
		Region:       pulumi.String(conf.Get("redisregion")),
		AutoScale:    pulumi.Bool(false),
		Eviction:     pulumi.Bool(true),
		Tls:          pulumi.Bool(true),
	})
	if err != nil {
		return err
	}

	ctx.Export("Redis Database Endpoint", redis.Endpoint)
	ctx.Export("Redis Database Post", redis.Port)
	ctx.Export("Redis Database Password", redis.Password)

	return nil
}

func setupSentryProject(ctx *pulumi.Context, project string, org string, team string) error {
	sentryProject, err := sentry.NewSentryProject(ctx, project, &sentry.SentryProjectArgs{
		Platform:     pulumi.String("go"),
		Organization: pulumi.String(org),
		Team:         pulumi.String(team),
		Name:         pulumi.String(project),
		Slug:         pulumi.String(project),
	})
	if err != nil {
		return err
	}

	ctx.Export("Sentry Project Name", sentryProject.ID())
	ctx.Export("Sentry Project ID", sentryProject.ProjectId)

	dsn := sentryProject.ID().ApplyT(func(name string) string {
		key, err := sentry.LookupSentryKey(ctx, &sentry.LookupSentryKeyArgs{
			First:        pulumi.BoolRef(true),
			Organization: org,
			Project:      project,
		})

		if err != nil {
			return ""
		}

		return key.DsnPublic
	})

	ctx.Export("Sentry Project DSN", dsn)
	return nil
}

func setupGCPProject(ctx *pulumi.Context, conf *config.Config) error {
	location := conf.Get("location")
	if location == "" {
		return fmt.Errorf("gcp location is not specified")
	}
	app, err := appengine.NewApplication(ctx, ctx.Project(), &appengine.ApplicationArgs{
		LocationId: pulumi.String(location),
	})
	if err != nil {
		return err
	}

	ctx.Export("app hostname", app.DefaultHostname)
	ctx.Export("app id", app.AppId)

	domain := conf.Get("domain")
	if domain != "" {
		mapping, err := appengine.NewDomainMapping(ctx, "domainMapping", &appengine.DomainMappingArgs{
			DomainName: pulumi.String(domain),
			SslSettings: &appengine.DomainMappingSslSettingsArgs{
				SslManagementType: pulumi.String("AUTOMATIC"),
			},
		})

		if err != nil {
			return err
		}

		ctx.Export("domain", mapping.DomainName)
	}
	return nil
}
