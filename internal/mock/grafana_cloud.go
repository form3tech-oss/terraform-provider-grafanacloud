package mock

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/naag/terraform-provider-grafanacloud/internal/api/grafana"
	"github.com/naag/terraform-provider-grafanacloud/internal/api/portal"
)

func (g *GrafanaCloud) createPortalAPIKey(w http.ResponseWriter, r *http.Request) {
	apiKey := &portal.APIKey{
		ID:    g.GetNextID(),
		Token: "very-secret",
	}
	fromJSON(apiKey, r)

	g.organisation.portalAPIKeys.AddKey(apiKey)
	sendResponse(w, apiKey, http.StatusCreated)
}

func (g *GrafanaCloud) listPortalAPIKeys(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, g.organisation.portalAPIKeys, http.StatusOK)
}

func (g *GrafanaCloud) deletePortalAPIKey(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	g.organisation.portalAPIKeys.DeleteByName(name)
	sendResponse(w, nil, http.StatusNoContent)
}

func (g *GrafanaCloud) createGrafanaAPIKeyProxy(w http.ResponseWriter, r *http.Request) {
	stackName := chi.URLParam(r, "stack")
	input := &portal.CreateGrafanaAPIKeyInput{}
	fromJSON(input, r)

	apiKey := &grafana.APIKey{
		Name: input.Name,
		Role: input.Role,
		ID:   g.GetNextID(),
		Key:  "very-secret",
	}

	if input.SecondsToLive > 0 {
		expiresAt := time.Now().Add(time.Duration(input.SecondsToLive) * time.Second)
		apiKey.Expiration = expiresAt.Format(time.RFC3339)
	}

	g.organisation.stackAPIKeys[stackName].AddKey(apiKey)
	sendResponse(w, apiKey, http.StatusCreated)
}

func (g *GrafanaCloud) listStacks(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, g.organisation.stackList, http.StatusOK)
}

func (g *GrafanaCloud) createStack(w http.ResponseWriter, r *http.Request) {
	stack := &portal.Stack{
		HmInstancePromID:  g.GetNextID(),
		HmInstancePromURL: "https://prometheus-instance",
		AmInstanceID:      g.GetNextID(),
	}
	fromJSON(stack, r)

	stack.ID = g.GetNextID()
	stack.OrgID = g.GetNextID()
	stack.OrgSlug = g.organisation.name
	stack.OrgName = g.organisation.name
	if stack.URL == "" {
		stack.URL = fmt.Sprintf("%s/grafana/%s", g.URL(), stack.Slug)
	}

	g.organisation.stackList.AddStack(stack)
	g.organisation.stackAPIKeys[stack.Slug] = &grafana.ListAPIKeysOutput{}
	sendResponse(w, stack, http.StatusCreated)
}

func (g *GrafanaCloud) deleteStack(w http.ResponseWriter, r *http.Request) {
	stackSlug := chi.URLParam(r, "stack")
	g.organisation.stackList.DeleteBySlug(stackSlug)
	delete(g.organisation.stackAPIKeys, stackSlug)
	sendResponse(w, nil, http.StatusNoContent)
}
