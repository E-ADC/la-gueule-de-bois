package handlers

import (
	"encoding/json"
	"net/http"

	"gueuledebois/backend/internal/services"
)

// DemandeAmiHandler expose UC21 (envoyer une demande d'ami) et sa réponse
// (accepter/refuser).
type DemandeAmiHandler struct {
	demandes *services.DemandeAmiService
}

func NewDemandeAmiHandler(demandes *services.DemandeAmiService) *DemandeAmiHandler {
	return &DemandeAmiHandler{demandes: demandes}
}

type envoyerDemandeAmiRequest struct {
	Pseudo string `json:"pseudo"`
}

// Envoyer — POST /api/amis/demandes (UC21).
func (h *DemandeAmiHandler) Envoyer(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	var req envoyerDemandeAmiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Pseudo == "" {
		mapAndWriteError(w, services.ErrValidation)
		return
	}

	if err := h.demandes.Envoyer(r.Context(), user.ID, req.Pseudo); err != nil {
		mapAndWriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// ListRecues — GET /api/amis/demandes (demandes en attente reçues).
func (h *DemandeAmiHandler) ListRecues(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	list, err := h.demandes.ListRecues(r.Context(), user.ID)
	if err != nil {
		mapAndWriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}

type repondreDemandeAmiRequest struct {
	// Action : "accepter" ou "refuser".
	Action string `json:"action"`
}

// Repondre — POST /api/amis/demandes/{id}/repondre.
func (h *DemandeAmiHandler) Repondre(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	demandeID, err := idFromURL(r, "id")
	if err != nil {
		mapAndWriteError(w, services.ErrValidation)
		return
	}

	var req repondreDemandeAmiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mapAndWriteError(w, services.ErrValidation)
		return
	}

	var accepter bool
	switch req.Action {
	case "accepter":
		accepter = true
	case "refuser":
		accepter = false
	default:
		mapAndWriteError(w, services.ErrValidation)
		return
	}

	if err := h.demandes.Repondre(r.Context(), user.ID, demandeID, accepter); err != nil {
		mapAndWriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListAmis — GET /api/amis : liste des amis (demandes acceptées).
func (h *DemandeAmiHandler) ListAmis(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	list, err := h.demandes.ListAmis(r.Context(), user.ID)
	if err != nil {
		mapAndWriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}
