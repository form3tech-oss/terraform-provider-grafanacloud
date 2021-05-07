package grafanacloud

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/form3tech-oss/terraform-provider-grafanacloud/internal/api/portal"
)

const (
	Name = "terraform-provider-grafanacloud"
	Addr = "github.com/form3tech-oss/grafanacloud"

	EnvURL            = "GRAFANA_CLOUD_URL"
	EnvOrganisation   = "GRAFANA_CLOUD_ORGANISATION"
	EnvAPIKey         = "GRAFANA_CLOUD_API_KEY"
	EnvTempKeyExpires = "GRAFANA_CLOUD_TEMP_KEY_EXPIRES"
	EnvTempKeyPrefix  = "GRAFANA_CLOUD_TEMP_KEY_PREFIX"
)

type Provider struct {
	Client       *portal.Client
	Organisation string
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
				"temp_key_expires": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: fmt.Sprintf("Time after which temporary Grafana API admin tokens used to read Grafana API resources expire. Might also be provided via `%s`", EnvTempKeyExpires),
					DefaultFunc: schema.EnvDefaultFunc(EnvTempKeyExpires, portal.TempKeyDefaultExpires),
				},
				"temp_key_prefix": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: fmt.Sprintf("Prefix for temporary Grafana API admin tokens used to read Grafana API resources. Might also be provided via `%s`", EnvTempKeyPrefix),
					DefaultFunc: schema.EnvDefaultFunc(EnvTempKeyPrefix, portal.TempKeyDefaultPrefix),
				},
			},
		}

		p.ConfigureContextFunc = ConfigureProvider(p, version)

		return p
	}
}

func ConfigureProvider(p *schema.Provider, version string) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		org := d.Get("organisation").(string)

		c, err := buildClient(p, d, version)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		err = c.AuthTest(ctx, org)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return &Provider{
			Client:       c,
			Organisation: org,
		}, nil
	}
}

func buildClient(p *schema.Provider, d *schema.ResourceData, version string) (*portal.Client, error) {
	url := d.Get("url").(string)
	apiKey := d.Get("api_key").(string)
	userAgent := p.UserAgent(Name, version)

	opts := []portal.ClientOpt{
		portal.WithUserAgent(userAgent),
	}

	if tempKeyExpires, ok := d.GetOk("temp_key_expires"); ok {
		d := time.Duration(tempKeyExpires.(int))
		opts = append(opts, portal.WithTempKeyExpires(d*time.Second))
	}

	if tempKeyPrefix, ok := d.GetOk("temp_key_prefix"); ok {
		opts = append(opts, portal.WithTempKeyPrefix(tempKeyPrefix.(string)))
	}

	return portal.NewClient(url, apiKey, opts...)
}
