package mock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/form3tech-oss/terraform-provider-grafanacloud/internal/api/grafana"
	"github.com/form3tech-oss/terraform-provider-grafanacloud/internal/api/portal"
)

type GrafanaCloud struct {
	organisation *organisation
	server       *httptest.Server
	nextID       int
}

type organisation struct {
	name          string
	stackList     *portal.ListStacksOutput
	portalAPIKeys *portal.ListAPIKeysOutput
	stackAPIKeys  map[string]*grafana.ListAPIKeysOutput
}

type errorResponse struct {
	Message string `json:"message"`
}

func (g *GrafanaCloud) Start() *GrafanaCloud {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)

	r.Post("/api/instances", g.createStack)
	r.Get("/api/orgs/{org}/instances", g.listStacks)
	r.Delete("/api/instances/{stack}", g.deleteStack)

	r.Post("/api/orgs/{org}/api-keys", g.createPortalAPIKey)
	r.Get("/api/orgs/{org}/api-keys", g.listPortalAPIKeys)
	r.Delete("/api/orgs/{org}/api-keys/{name}", g.deletePortalAPIKey)

	r.Post("/api/instances/{stack}/api/auth/keys", g.createGrafanaAPIKeyProxy)

	// Grafana Cloud API doesn't really offer routes at /api/grafana. These are just provided
	// here so that we can mock the Grafana API running inside Grafana Cloud stacks.
	r.Get("/api/grafana/{stack}/api/auth/keys", g.listGrafanaAPIKeys)
	r.Delete("/api/grafana/{stack}/api/auth/keys/{id}", g.deleteGrafanaAPIKey)

	g.server = httptest.NewServer(r)
	return g
}

func NewGrafanaCloud(org string) *GrafanaCloud {
	return &GrafanaCloud{
		organisation: &organisation{
			name:          org,
			stackList:     &portal.ListStacksOutput{},
			portalAPIKeys: &portal.ListAPIKeysOutput{},
			stackAPIKeys:  make(map[string]*grafana.ListAPIKeysOutput),
		},
	}
}

func (g *GrafanaCloud) Close() {
	g.server.Close()
}

func (g *GrafanaCloud) URL() string {
	return fmt.Sprintf("%s/api", g.server.URL)
}

func fromJSON(d interface{}, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, d); err != nil {
		panic(err)
	}
}

func sendResponse(w http.ResponseWriter, v interface{}, status int) {
	if v != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	}

	w.WriteHeader(status)

	if v != nil {
		if err := json.NewEncoder(w).Encode(v); err != nil {
			panic(err)
		}
	}
}

func sendError(w http.ResponseWriter, err error) {
	resp := &errorResponse{
		Message: err.Error(),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}

func (g *GrafanaCloud) GetNextID() int {
	g.nextID += 1
	return g.nextID
}
