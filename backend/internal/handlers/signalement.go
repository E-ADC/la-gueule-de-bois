package handlers

import (
	"encoding/json"
	"net/http"

	"gueuledebois/backend/internal/models"
	"gueuledebois/backend/internal/services"
)

// SignalementHandler expose UC13 (signaler un témoignage) et UC22 (traiter
// un signalement, acteur Modérateur).
type SignalementHandler struct {
	signalements *services.SignalementService
	temoignages  *services.TemoignageService
}

func NewSignalementHandler(signalements *services.SignalementService, temoignages *services.TemoignageService) *SignalementHandler {
	return &SignalementHandler{signalements: signalements, temoignages: temoignages}
}

type signalerRequest struct {
	Motif string `json:"motif"`
}

// Report — POST /api/temoignages/{id}/signalements (UC13).
func (h *SignalementHandler) Report(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	temoignageID, err := idFromURL(r, "id")
	if err != nil {
		mapAndWriteError(w, services.ErrValidation)
		return
	}

	var req signalerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mapAndWriteError(w, services.ErrValidation)
		return
	}

	if err := h.signalements.Report(r.Context(), user.ID, temoignageID, req.Motif); err != nil {
		mapAndWriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// signalementView enrichit un signalement avec le contenu du témoignage
// signalé, pour affichage direct sur la page de modération (UC22) sans
// aller-retour supplémentaire côté frontend.
type signalementView struct {
	models.Signalement
	TemoignageContenu string `json:"temoignageContenu"`
}

// ListEnAttente — GET /api/signalements (UC22 étape 1, modérateur uniquement).
func (h *SignalementHandler) ListEnAttente(w http.ResponseWriter, r *http.Request) {
	list, err := h.signalements.ListEnAttente(r.Context())
	if err != nil {
		mapAndWriteError(w, err)
		return
	}

	views := make([]signalementView, 0, len(list))
	for _, s := range list {
		view := signalementView{Signalement: s}
		if s.TemoignageID != nil {
			if t, err := h.temoignages.Get(r.Context(), *s.TemoignageID); err == nil {
				view.TemoignageContenu = t.Contenu
			}
		}
		views = append(views, view)
	}
	writeJSON(w, http.StatusOK, views)
}

type traiterRequest struct {
	// Action : "rejeter" (le signalement est clos, le témoignage conservé)
	// ou "supprimer" (le témoignage est supprimé).
	Action string `json:"action"`
}

// Traiter — POST /api/signalements/{id}/traiter (UC22 étapes 2-4, modérateur uniquement).
func (h *SignalementHandler) Traiter(w http.ResponseWriter, r *http.Request) {
	moderateur := UserFromContext(r.Context())
	signalementID, err := idFromURL(r, "id")
	if err != nil {
		mapAndWriteError(w, services.ErrValidation)
		return
	}

	var req traiterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mapAndWriteError(w, services.ErrValidation)
		return
	}

	var supprimer bool
	switch req.Action {
	case "supprimer":
		supprimer = true
	case "rejeter":
		supprimer = false
	default:
		mapAndWriteError(w, services.ErrValidation)
		return
	}

	if err := h.signalements.Traiter(r.Context(), moderateur.ID, signalementID, supprimer); err != nil {
		mapAndWriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
