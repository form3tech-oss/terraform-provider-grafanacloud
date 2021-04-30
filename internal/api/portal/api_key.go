package portal

import (
	"fmt"

	"github.com/naag/terraform-provider-grafanacloud/internal/api/grafana"
	"github.com/naag/terraform-provider-grafanacloud/internal/util"
)

type CreateAPIKeyInput struct {
	Name         string `json:"name"`
	Role         string `json:"role"`
	Organisation string `json:"-"`
}

type CreateGrafanaAPIKeyInput struct {
	Name          string `json:"name"`
	Role          string `json:"role"`
	SecondsToLive int    `json:"secondsToLive"`
	Stack         string `json:"-"`
}

type CreateGrafanaAPIKeyOutput grafana.APIKey

type ListAPIKeysOutput struct {
	Items []*APIKey
}

type APIKey struct {
	ID         int
	Name       string
	Role       string
	Token      string
	Expiration string
}

// This function creates a API key inside the Grafana instance running in stack `stack`. It's used in order
// to provision API keys inside Grafana while just having access to a Grafana Cloud API key.
//
// Plese note that this is a beta feature and might change in the future.
//
// See https://grafana.com/docs/grafana-cloud/api/#create-grafana-api-keys for more information.
func (c *Client) CreateGrafanaAPIKey(r *CreateGrafanaAPIKeyInput) (*CreateGrafanaAPIKeyOutput, error) {
	url := fmt.Sprintf("instances/%s/api/auth/keys", r.Stack)
	resp, err := c.client.R().
		SetBody(r).
		SetResult(&CreateGrafanaAPIKeyOutput{}).
		Post(url)

	if err := util.HandleError(err, resp, "Failed to create Grafana API key through Grafana Cloud proxy route"); err != nil {
		return nil, err
	}

	return resp.Result().(*CreateGrafanaAPIKeyOutput), nil
}

func (c *Client) CreateAPIKey(r *CreateAPIKeyInput) (*APIKey, error) {
	url := fmt.Sprintf("orgs/%s/api-keys", r.Organisation)
	resp, err := c.client.R().
		SetBody(r).
		SetResult(&APIKey{}).
		Post(url)

	if err := util.HandleError(err, resp, "failed to create Grafana Cloud Portal API key"); err != nil {
		return nil, err
	}

	return resp.Result().(*APIKey), nil
}

func (c *Client) ListAPIKeys(org string) (*ListAPIKeysOutput, error) {
	url := fmt.Sprintf("orgs/%s/api-keys", org)
	resp, err := c.client.R().
		SetResult(&ListAPIKeysOutput{}).
		Get(url)

	if err := util.HandleError(err, resp, "failed to read Grafana Cloud Portal API key"); err != nil {
		return nil, err
	}

	return resp.Result().(*ListAPIKeysOutput), nil
}

func (c *Client) DeleteAPIKey(org string, keyName string) error {
	url := fmt.Sprintf("orgs/%s/api-keys/%s", org, keyName)
	resp, err := c.client.R().
		Delete(url)

	if err := util.HandleError(err, resp, "failed to delete Grafana Cloud Portal API key"); err != nil {
		return err
	}

	return nil
}

func (l *ListAPIKeysOutput) AddKey(k *APIKey) {
	l.Items = append(l.Items, k)
}

func (l *ListAPIKeysOutput) FindByName(name string) *APIKey {
	for _, k := range l.Items {
		if k.Name == name {
			return k
		}
	}

	return nil
}

func (l *ListAPIKeysOutput) DeleteByName(name string) {
	newKeys := make([]*APIKey, 0)

	for _, k := range l.Items {
		if k.Name != name {
			newKeys = append(newKeys, k)
		}
	}

	l.Items = newKeys
}
