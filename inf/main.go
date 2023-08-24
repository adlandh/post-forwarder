package main

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/appengine"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumiverse/pulumi-sentry/sdk/go/sentry"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		location, ok := ctx.GetConfig(ctx.Project() + ":location")
		if !ok {
			return fmt.Errorf("gcp location is not specified")
		}
		app, err := appengine.NewApplication(ctx, ctx.Project(), &appengine.ApplicationArgs{
			LocationId: pulumi.String(location),
		})
		if err != nil {
			return err
		}

		ctx.Export("app hostname", app.DefaultHostname)

		domain, ok := ctx.GetConfig(ctx.Project() + ":domain")

		if ok {
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

		project, okProject := ctx.GetConfig(ctx.Project() + ":sentry-project")
		org, okOrg := ctx.GetConfig(ctx.Project() + ":sentry-org")
		team, okTeam := ctx.GetConfig(ctx.Project() + ":sentry-team")
		if okProject && okOrg && okTeam {
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
