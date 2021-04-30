package grafanacloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceStack() *schema.Resource {
	s := baseStackSchema()
	s["slug"].Required = true
	s["slug"].Computed = false

	return &schema.Resource{
		Description: "Reads a single Grafana Cloud stack from the organisation by the given name.",
		ReadContext: dataSourceStackRead,
		Schema:      s,
	}
}
func dataSourceStackRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	p := m.(*Provider)
	slug := d.Get("slug").(string)

	stackList, err := listStacks(p)
	if err != nil {
		return diag.FromErr(err)
	}

	stack := stackList.FindBySlug(slug)
	if stack == nil {
		return diag.Errorf("Couldn't find stack with slug `%s`", slug)
	}

	d.SetId(fmt.Sprint(stack.ID))

	if err := d.Set("name", stack.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("slug", stack.Slug); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("prometheus_url", stack.HmInstancePromURL); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("prometheus_user_id", stack.HmInstancePromID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("alertmanager_url", stack.AmInstanceURL); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("alertmanager_user_id", stack.AmInstanceID); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
