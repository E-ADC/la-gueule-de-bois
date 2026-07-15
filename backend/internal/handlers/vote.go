package handlers

import (
	"encoding/json"
	"net/http"

	"gueuledebois/backend/internal/models"
	"gueuledebois/backend/internal/services"
)

// VoteHandler expose UC12 (swiper/voter sur un témoignage).
type VoteHandler struct {
	votes *services.VoteService
}

func NewVoteHandler(votes *services.VoteService) *VoteHandler {
	return &VoteHandler{votes: votes}
}

type voteRequest struct {
	// Valeur : 1 (positif) ou -1 (négatif) — cf. UC12 "swipe positif ou négatif".
	Valeur int `json:"valeur"`
}

// Cast — POST /api/temoignages/{id}/votes (UC12).
func (h *VoteHandler) Cast(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	temoignageID, err := idFromURL(r, "id")
	if err != nil {
		mapAndWriteError(w, services.ErrValidation)
		return
	}

	var req voteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || (req.Valeur != 1 && req.Valeur != -1) {
		mapAndWriteError(w, services.ErrValidation)
		return
	}

	v, err := h.votes.Cast(r.Context(), user.ID, temoignageID, models.VoteValeur(req.Valeur))
	if err != nil {
		mapAndWriteError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, v)
}
