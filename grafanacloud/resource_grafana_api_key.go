package grafanacloud

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/naag/terraform-provider-grafanacloud/internal/api/portal"
)

var (
	grafanaApiKeyRoles = []string{"Viewer", "Editor", "Admin"}
)

func resourceGrafanaApiKey() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a single API key on a Grafana instance inside a Grafana Cloud stack. Notice that the key value will be stored in Terraform state, so make sure to manage your Terraform state safely (see https://www.terraform.io/docs/language/state/sensitive-data.html).",
		CreateContext: resourceApiKeyCreate,
		ReadContext:   resourceApiKeyRead,
		DeleteContext: resourceApiKeyDelete,
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
			"stack": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Grafana Cloud stack to create this API key in.",
			},
			"role": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  fmt.Sprintf("Role of the API key. Might be one of %s. See https://grafana.com/docs/grafana-cloud/api/#create-api-key for details.", grafanaApiKeyRoles),
				ValidateFunc: ValidateGrafanaApiKeyRole(),
			},
			"seconds_to_live": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Time in seconds after which the API key automatically expires",
			},
			"expiration": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time at which the API key expires (ISO8601 format). Blank if the API key never expires.",
			},
			"is_expired": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Whether or not the API key has expired. This field is used internally in order to recreate expired API keys. Set this to `true` to not recreate expired API keys.",
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

func ValidateGrafanaApiKeyRole() schema.SchemaValidateFunc {
	return validation.StringInSlice(grafanaApiKeyRoles, false)
}

func resourceApiKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	p := m.(*Provider)

	req := &portal.CreateGrafanaAPIKeyInput{
		Name:          d.Get("name").(string),
		Role:          d.Get("role").(string),
		Stack:         d.Get("stack").(string),
		SecondsToLive: d.Get("seconds_to_live").(int),
	}

	resp, err := p.Client.CreateGrafanaAPIKey(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("key", resp.Key)
	d.SetId(strconv.Itoa(resp.ID))

	return resourceApiKeyRead(ctx, d, m)
}

func resourceApiKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	p := m.(*Provider)

	stack := d.Get("stack").(string)
	client, cleanup, err := p.Client.GetAuthedGrafanaClient(ctx, p.Organisation, stack)
	if err != nil {
		return diag.FromErr(err)
	}

	if cleanup != nil {
		defer cleanup()
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	apiKeys, err := client.ListAPIKeys(ctx, true)
	if err != nil {
		return diag.FromErr(err)
	}

	apiKey := apiKeys.FindByID(id)
	if apiKey == nil {
		d.SetId("")
		return diags
	}

	expired, err := apiKey.IsExpired()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", apiKey.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("role", apiKey.Role); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("expiration", apiKey.Expiration); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("is_expired", expired); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceApiKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	p := m.(*Provider)

	stack := d.Get("stack").(string)
	client, cleanup, err := p.Client.GetAuthedGrafanaClient(ctx, p.Organisation, stack)
	if err != nil {
		return diag.FromErr(err)
	}

	if cleanup != nil {
		defer cleanup()
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteAPIKey(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}
