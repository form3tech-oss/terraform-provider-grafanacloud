package portal

import (
	"fmt"
	"log"
	"time"

	"github.com/naag/terraform-provider-grafanacloud/internal/api/grafana"
)

// The Grafana Cloud API is disconnected from the Grafana API on the stacks unfortunately. That's why we can't use
// the Grafana Cloud API key to fully manage API keys on the Grafana API. The only thing we can do is to create
// a temporary Admin key, and create a Grafana API client with that.
func (c *Client) GetAuthedGrafanaClient(orgName, stackName string) (*grafana.Client, func() error, error) {
	stack, err := c.GetStack(orgName, stackName)
	if err != nil {
		return nil, nil, err
	}

	if stack == nil {
		return nil, nil, fmt.Errorf("failed to find stack by name %s", stackName)
	}

	name := fmt.Sprintf("terraform-provider-grafanacloud-tmp-%d", time.Now().UnixNano())
	req := &CreateGrafanaAPIKeyInput{
		Name:          name,
		Role:          "Admin",
		SecondsToLive: 60,
		Stack:         stackName,
	}

	apiKey, err := c.CreateGrafanaAPIKey(req)
	if err != nil {
		return nil, nil, err
	}

	log.Printf("[DEBUG] created a temporary admin API key `%s` on Grafana stack `%s`", apiKey.Name, stack.Slug)

	client, err := grafana.NewClient(stack.URL, apiKey.Key, grafana.WithUserAgent(c.client.Header.Get("User-Agent")))
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() error {
		err = client.DeleteAPIKey(apiKey.ID)
		if err != nil {
			log.Printf("[ERROR] failed deleting temporary admin API key `%s` on Grafana stack `%s`", apiKey.Name, stack.Slug)
			return err
		}

		log.Printf("[DEBUG] deleted temporary admin API key `%s` on Grafana stack `%s`", apiKey.Name, stack.Slug)
		return nil
	}

	return client, cleanup, nil
}
