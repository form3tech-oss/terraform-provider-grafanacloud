package portal

import (
	"context"
	"fmt"

	"github.com/naag/terraform-provider-grafanacloud/internal/util"
)

type CreateStackInput struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
	URL  string `json:"url"`
}

type ListStacksOutput struct {
	Items []*Stack
}

type Stack struct {
	ID                   int
	OrgID                int
	OrgSlug              string
	OrgName              string
	Name                 string
	URL                  string
	Status               string
	Slug                 string
	HmInstancePromID     int
	HmInstancePromURL    string
	HmInstancePromStatus string
	AmInstanceID         int
	AmInstanceURL        string
}

func (c *Client) CreateStack(ctx context.Context, r *CreateStackInput) (*Stack, error) {
	url := "instances"
	resp, err := c.client.R().
		SetBody(r).
		SetResult(&Stack{}).
		SetContext(ctx).
		Post(url)

	if err := util.HandleError(err, resp, "failed to create Grafana Cloud stack"); err != nil {
		return nil, err
	}

	return resp.Result().(*Stack), nil
}

func (c *Client) ListStacks(ctx context.Context, org string) (*ListStacksOutput, error) {
	url := fmt.Sprintf("orgs/%s/instances", org)
	resp, err := c.client.R().
		SetResult(&ListStacksOutput{}).
		SetContext(ctx).
		Get(url)

	if err := util.HandleError(err, resp, "failed to list Grafana Cloud stacks"); err != nil {
		return nil, err
	}

	return resp.Result().(*ListStacksOutput), nil
}

func (c *Client) GetStack(ctx context.Context, org, stackSlug string) (*Stack, error) {
	stacks, err := c.ListStacks(ctx, org)
	if err != nil {
		return nil, err
	}

	stack := stacks.FindBySlug(stackSlug)
	return stack, nil
}

func (c *Client) DeleteStack(ctx context.Context, stackSlug string) error {
	url := fmt.Sprintf("instances/%s", stackSlug)
	resp, err := c.client.R().
		SetContext(ctx).
		Delete(url)

	if err := util.HandleError(err, resp, "failed to delete Grafana Cloud stack"); err != nil {
		return err
	}

	return nil
}

func (l *ListStacksOutput) AddStack(s *Stack) {
	l.Items = append(l.Items, s)
}

func (l *ListStacksOutput) FindBySlug(slug string) *Stack {
	for _, stack := range l.Items {
		if stack.Slug == slug {
			return stack
		}
	}

	return nil
}

func (l *ListStacksOutput) DeleteBySlug(slug string) {
	newItems := make([]*Stack, 0)

	for _, stack := range l.Items {
		if stack.Slug != slug {
			newItems = append(newItems, stack)
		}
	}

	l.Items = newItems
}
