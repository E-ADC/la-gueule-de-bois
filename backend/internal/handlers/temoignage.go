package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"gueuledebois/backend/internal/models"
	"gueuledebois/backend/internal/services"
)

// TemoignageHandler expose UC11 (ajouter un témoignage) et UC09 (inviter
// un témoin, pré-condition de UC11).
type TemoignageHandler struct {
	temoignages *services.TemoignageService
	votes       *services.VoteService
	profile     *services.ProfileService
}

func NewTemoignageHandler(temoignages *services.TemoignageService, votes *services.VoteService, profile *services.ProfileService) *TemoignageHandler {
	return &TemoignageHandler{temoignages: temoignages, votes: votes, profile: profile}
}

type inviteTemoinRequest struct {
	Pseudo string `json:"pseudo"`
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Pseudo) == "" {
		mapAndWriteError(w, services.ErrValidation)
		return
	}

	if err := h.temoignages.InviteTemoin(r.Context(), user.ID, soireeID, req.Pseudo); err != nil {
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

// temoignageView enrichit un témoignage du pseudo de son auteur, des
// compteurs de votes et du vote de l'utilisateur courant, pour un affichage
// complet sans aller-retours supplémentaires côté frontend.
type temoignageView struct {
	models.Temoignage
	AuteurPseudo  string `json:"auteurPseudo"`
	VotesPositifs int    `json:"votesPositifs"`
	VotesNegatifs int    `json:"votesNegatifs"`
	MonVote       *int   `json:"monVote"`
}

// ListInvitedSoirees — GET /api/soirees/invitations : les soirées où
// l'utilisateur connecté a été invité comme témoin (UC09).
func (h *TemoignageHandler) ListInvitedSoirees(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	list, err := h.temoignages.ListInvitedSoirees(r.Context(), user.ID)
	if err != nil {
		mapAndWriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// ListBySoiree — GET /api/soirees/{id}/temoignages.
func (h *TemoignageHandler) ListBySoiree(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
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

	views := make([]temoignageView, 0, len(list))
	for _, t := range list {
		view := temoignageView{Temoignage: t}
		if auteur, err := h.profile.GetPublicProfile(r.Context(), t.AuteurID); err == nil {
			view.AuteurPseudo = auteur.Pseudo
		}
		positifs, negatifs, err := h.votes.Counts(r.Context(), t.ID)
		if err == nil {
			view.VotesPositifs = positifs
			view.VotesNegatifs = negatifs
		}
		if monVote, err := h.votes.MonVote(r.Context(), user.ID, t.ID); err == nil && monVote != nil {
			valeur := int(monVote.Valeur)
			view.MonVote = &valeur
		}
		views = append(views, view)
	}
	writeJSON(w, http.StatusOK, views)
}
