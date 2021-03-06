package grafana

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/form3tech-oss/terraform-provider-grafanacloud/internal/util"
	"github.com/relvacode/iso8601"
)

type ListAPIKeysOutput struct {
	Keys []*APIKey
}

type APIKey struct {
	ID         int
	Name       string
	Key        string
	Expiration string
	Role       string
}

func (c *Client) ListAPIKeys(ctx context.Context, includeExpired bool) (*ListAPIKeysOutput, error) {
	var apiKeys []*APIKey
	url := "api/auth/keys"

	resp, err := c.client.R().
		SetResult(&apiKeys).
		SetQueryParam("includeExpired", strconv.FormatBool(includeExpired)).
		SetContext(ctx).
		Get(url)

	if err := util.HandleError(err, resp, "failed to list Grafana API keys"); err != nil {
		return nil, err
	}

	return &ListAPIKeysOutput{
		Keys: apiKeys,
	}, nil
}

func (c *Client) DeleteAPIKey(ctx context.Context, id int) error {
	url := fmt.Sprintf("api/auth/keys/%d", id)

	resp, err := c.client.R().
		SetContext(ctx).
		Delete(url)

	if err := util.HandleError(err, resp, "failed to delete Grafana API key"); err != nil {
		return err
	}

	return nil
}

func (k *APIKey) IsExpired() (bool, error) {
	if k.Expiration == "" {
		return false, nil
	}

	expires, err := iso8601.ParseString(k.Expiration)
	if err != nil {
		return false, err
	}

	now := time.Now()
	return now.After(expires), nil
}

func (l *ListAPIKeysOutput) AddKey(k *APIKey) {
	l.Keys = append(l.Keys, k)
}

func (l *ListAPIKeysOutput) FindByID(id int) *APIKey {
	for _, k := range l.Keys {
		if k.ID == id {
			return k
		}
	}

	return nil
}

func (l *ListAPIKeysOutput) DeleteByID(id int) {
	newKeys := make([]*APIKey, 0)

	for _, k := range l.Keys {
		if k.ID != id {
			newKeys = append(newKeys, k)
		}
	}

	l.Keys = newKeys
}
