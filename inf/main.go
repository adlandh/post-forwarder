package main

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/appengine"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"github.com/pulumiverse/pulumi-sentry/sdk/go/sentry"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		conf := config.New(ctx, "")
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

		project := conf.Get("sentry-project")
		org := conf.Get("sentry-org")
		team := conf.Get("sentry-team")
		if project != "" && org != "" && team != "" {
			sentryProject, err := sentry.NewSentryProject(ctx, project, &sentry.SentryProjectArgs{
				Platform:     pulumi.String("go"),
				Organization: pulumi.String(org),
				Team:         pulumi.String(team),
			})
			if err != nil {
				return err
			}

			ctx.Export("Sentry Project Name", sentryProject.ID())
			ctx.Export("Sentry Project ID", sentryProject.ProjectId)

			key, err := sentry.NewSentryKey(ctx, "default", &sentry.SentryKeyArgs{
				Organization: pulumi.String(org),
				Project:      sentryProject.ID(),
			})
			if err != nil {
				return err
			}

			ctx.Export("Sentry Project DSN", key.DsnPublic)
		}

		return nil
	})
}
