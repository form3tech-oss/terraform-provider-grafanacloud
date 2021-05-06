package portal

import (
	"context"
	"fmt"

	"github.com/naag/terraform-provider-grafanacloud/internal/util"
)

func (c *Client) AuthTest(ctx context.Context, org string) error {
	url := fmt.Sprintf("orgs/%s/instances", org)
	resp, err := c.client.R().
		SetContext(ctx).
		Get(url)

	if err := util.HandleError(err, resp, "failed to test connection with Grafana Cloud API. Please check API key and organisation"); err != nil {
		return err
	}

	return nil
}
