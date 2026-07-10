package handlers

import (
	"encoding/json"
	"net/http"

	"gueuledebois/backend/internal/services"
)

// TemoignageHandler expose UC11 (ajouter un témoignage) et UC09 (inviter
// un témoin, pré-condition de UC11).
type TemoignageHandler struct {
	temoignages *services.TemoignageService
}

func NewTemoignageHandler(temoignages *services.TemoignageService) *TemoignageHandler {
	return &TemoignageHandler{temoignages: temoignages}
}

type inviteTemoinRequest struct {
	InviteID int64 `json:"inviteId"`
}

// InviteTemoin — POST /api/soirees/{id}/temoins (UC09).
func (h *TemoignageHandler) InviteTemoin(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	soireeID, err := idFromURL(r, "id")
	if err != nil {
		mapAndWriteError(w, services.ErrValidation)
		return
	}

	var req inviteTemoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.InviteID <= 0 {
		mapAndWriteError(w, services.ErrValidation)
		return
	}

	if err := h.temoignages.InviteTemoin(r.Context(), user.ID, soireeID, req.InviteID); err != nil {
		mapAndWriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

type temoignageRequest struct {
	Contenu string `json:"contenu"`
}

// Add — POST /api/soirees/{id}/temoignages (UC11, inclut UC16).
func (h *TemoignageHandler) Add(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	soireeID, err := idFromURL(r, "id")
	if err != nil {
		mapAndWriteError(w, services.ErrValidation)
		return
	}

	var req temoignageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mapAndWriteError(w, services.ErrValidation)
		return
	}

	t, err := h.temoignages.Add(r.Context(), user.ID, soireeID, req.Contenu)
	if err != nil {
		mapAndWriteError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, t)
}

// ListBySoiree — GET /api/soirees/{id}/temoignages.
func (h *TemoignageHandler) ListBySoiree(w http.ResponseWriter, r *http.Request) {
	soireeID, err := idFromURL(r, "id")
	if err != nil {
		mapAndWriteError(w, services.ErrValidation)
		return
	}

	list, err := h.temoignages.ListBySoiree(r.Context(), soireeID)
	if err != nil {
		mapAndWriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}
