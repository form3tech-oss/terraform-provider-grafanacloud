package grafanacloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/naag/terraform-provider-grafanacloud/internal/api/portal"
)

const (
	Name = "terraform-provider-grafanacloud"
	Addr = "github.com/form3tech-oss/grafanacloud"

	EnvURL          = "GRAFANA_CLOUD_URL"
	EnvOrganisation = "GRAFANA_CLOUD_ORGANISATION"
	EnvAPIKey       = "GRAFANA_CLOUD_API_KEY"
)

type Provider struct {
	Client       *portal.Client
	Organisation string
	UserAgent    string
}

func NewProvider(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			ResourcesMap: map[string]*schema.Resource{
				"grafanacloud_stack":           resourceStack(),
				"grafanacloud_grafana_api_key": resourceGrafanaApiKey(),
				"grafanacloud_portal_api_key":  resourcePortalApiKey(),
			},
			DataSourcesMap: map[string]*schema.Resource{
				"grafanacloud_stacks": dataSourceStacks(),
				"grafanacloud_stack":  dataSourceStack(),
			},
			Schema: map[string]*schema.Schema{
				"url": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: fmt.Sprintf("Grafana Cloud API endpoint including the final `/api`. Might also be provided via `%s`.", EnvURL),
					DefaultFunc: schema.EnvDefaultFunc(EnvURL, "https://grafana.com/api"),
				},
				"api_key": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					Description: fmt.Sprintf("API key used to authenticate with the API. Must have `Admin` role if API keys need to be managed. Might also be provided via `%s`.", EnvAPIKey),
					DefaultFunc: schema.EnvDefaultFunc(EnvAPIKey, ""),
				},
				"organisation": {
					Type:        schema.TypeString,
					Required:    true,
					Description: fmt.Sprintf("Organisation which the API key belongs to (as slug name). Might also be provided via `%s`", EnvOrganisation),
					DefaultFunc: schema.EnvDefaultFunc(EnvOrganisation, ""),
				},
			},
		}

		p.ConfigureContextFunc = ConfigureProvider(version, p)

		return p
	}
}

func ConfigureProvider(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		apiKey := d.Get("api_key").(string)
		url := d.Get("url").(string)
		org := d.Get("organisation").(string)

		userAgent := p.UserAgent(Name, version)
		c, err := portal.NewClient(url, apiKey, portal.WithUserAgent(userAgent))
		if err != nil {
			return nil, diag.FromErr(err)
		}

		err = c.AuthTest(org)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return &Provider{
			Client:       c,
			Organisation: org,
			UserAgent:    userAgent,
		}, nil
	}
}
