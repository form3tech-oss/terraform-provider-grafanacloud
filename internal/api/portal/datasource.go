package portal

import (
	"context"
	"fmt"

	"github.com/form3tech-oss/terraform-provider-grafanacloud/internal/util"
)

type ListDatasourcesOutput struct {
	Items []*Datasource
}
type Datasource struct {
	ID            int
	InstanceID    int
	InstanceSlug  string
	Name          string
	Type          string
	URL           string
	BasicAuth     int
	BasicAuthUser string
}

func (c *Client) ListDatasources(ctx context.Context, stack string) (*ListDatasourcesOutput, error) {
	url := fmt.Sprintf("instances/%s/datasources", stack)
	resp, err := c.client.R().
		SetResult(&ListDatasourcesOutput{}).
		SetContext(ctx).
		Get(url)

	if err := util.HandleError(err, resp, "failed to list Grafana data sources"); err != nil {
		return nil, err
	}

	return resp.Result().(*ListDatasourcesOutput), nil
}

func (ds *Datasource) IsAlertmanager() bool {
	return ds.Type == "grafana-alertmanager-datasource"
}
