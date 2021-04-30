package portal

import (
	"fmt"

	"github.com/naag/terraform-provider-grafanacloud/internal/util"
)

func (c *Client) AuthTest(org string) error {
	url := fmt.Sprintf("orgs/%s/instances", org)
	resp, err := c.client.R().
		Get(url)

	if err := util.HandleError(err, resp, "failed to test connection with Grafana Cloud API. Please check API key and organisation"); err != nil {
		return err
	}

	return nil
}
