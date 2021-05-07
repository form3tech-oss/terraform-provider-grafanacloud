package mock

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (g *GrafanaCloud) listGrafanaAPIKeys(w http.ResponseWriter, r *http.Request) {
	stackName := chi.URLParam(r, "stack")
	sendResponse(w, g.organisation.stackAPIKeys[stackName].Keys, http.StatusOK)
}

func (g *GrafanaCloud) deleteGrafanaAPIKey(w http.ResponseWriter, r *http.Request) {
	stackName := chi.URLParam(r, "stack")
	keyID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		sendError(w, err)
		return
	}

	g.organisation.stackAPIKeys[stackName].DeleteByID(keyID)
	sendResponse(w, nil, http.StatusNoContent)
}
