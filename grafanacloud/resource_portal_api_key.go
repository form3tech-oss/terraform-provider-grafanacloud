package grafanacloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/naag/terraform-provider-grafanacloud/internal/api/portal"
)

var (
	portalApiKeyRoles = []string{"Viewer", "Editor", "Admin", "MetricsPublisher", "PluginPublisher"}
)

func resourcePortalApiKey() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a single API key on the Grafana Cloud portal (on the organisation level). Notice that the key value will be stored in Terraform state, so make sure to manage your Terraform state safely (see https://www.terraform.io/docs/language/state/sensitive-data.html).",
		CreateContext: resourcePortalApiKeyCreate,
		ReadContext:   resourcePortalApiKeyRead,
		DeleteContext: resourcePortalApiKeyDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the API key.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the API key.",
			},
			"role": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  fmt.Sprintf("Role of the API key. Might be one of %s. See https://grafana.com/docs/grafana-cloud/api/#create-api-key for details.", portalApiKeyRoles),
				ValidateFunc: ValidatePortalApiKeyRole(),
			},
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The generated API key.",
			},
		},
	}
}

func ValidatePortalApiKeyRole() schema.SchemaValidateFunc {
	return validation.StringInSlice(portalApiKeyRoles, false)
}

func resourcePortalApiKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	p := m.(*Provider)

	req := &portal.CreateAPIKeyInput{
		Name:         d.Get("name").(string),
		Role:         d.Get("role").(string),
		Organisation: p.Organisation,
	}

	resp, err := p.Client.CreateAPIKey(req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("key", resp.Token)
	d.SetId(resp.Name)

	return resourcePortalApiKeyRead(ctx, d, m)
}

func resourcePortalApiKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	p := m.(*Provider)

	resp, err := p.Client.ListAPIKeys(p.Organisation)
	if err != nil {
		return diag.FromErr(err)
	}

	portalKey := resp.FindByName(d.Id())
	if portalKey == nil {
		d.SetId("")
		return diags
	}

	if err := d.Set("name", portalKey.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("role", portalKey.Role); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourcePortalApiKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	p := m.(*Provider)

	err := p.Client.DeleteAPIKey(p.Organisation, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}
