package handlers

import (
	"encoding/json"
	"net/http"

	"gueuledebois/backend/internal/services"
)

// GroupeHandler expose UC18 (créer un groupe) et UC19 (rejoindre un groupe).
type GroupeHandler struct {
	groupes *services.GroupeService
}

func NewGroupeHandler(groupes *services.GroupeService) *GroupeHandler {
	return &GroupeHandler{groupes: groupes}
}

type groupeRequest struct {
	Nom string `json:"nom"`
}

// Create — POST /api/groupes (UC18).
func (h *GroupeHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	var req groupeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mapAndWriteError(w, services.ErrValidation)
		return
	}

	g, err := h.groupes.Create(r.Context(), user.ID, req.Nom)
	if err != nil {
		mapAndWriteError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, g)
}

// Join — POST /api/groupes/{id}/membres (UC19).
func (h *GroupeHandler) Join(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	groupeID, err := idFromURL(r, "id")
	if err != nil {
		mapAndWriteError(w, services.ErrValidation)
		return
	}

	if err := h.groupes.Join(r.Context(), user.ID, groupeID); err != nil {
		mapAndWriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListMine — GET /api/groupes (groupes dont l'utilisateur est membre).
func (h *GroupeHandler) ListMine(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	list, err := h.groupes.ListMine(r.Context(), user.ID)
	if err != nil {
		mapAndWriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// Get — GET /api/groupes/{id} (page de détail d'un groupe).
func (h *GroupeHandler) Get(w http.ResponseWriter, r *http.Request) {
	groupeID, err := idFromURL(r, "id")
	if err != nil {
		mapAndWriteError(w, services.ErrValidation)
		return
	}
	g, err := h.groupes.Get(r.Context(), groupeID)
	if err != nil {
		mapAndWriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, g)
}
