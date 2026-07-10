package handlers

import (
	"net/http"

	"gueuledebois/backend/internal/services"
)

// ProfileHandler expose UC05 (profil public), UC15 (badges), UC17/UC20
// (classements).
type ProfileHandler struct {
	profile *services.ProfileService
}

func NewProfileHandler(profile *services.ProfileService) *ProfileHandler {
	return &ProfileHandler{profile: profile}
}

// GetPublicProfile — GET /api/users/{id} (UC05).
func (h *ProfileHandler) GetPublicProfile(w http.ResponseWriter, r *http.Request) {
	id, err := idFromURL(r, "id")
	if err != nil {
		mapAndWriteError(w, services.ErrValidation)
		return
	}
	user, err := h.profile.GetPublicProfile(r.Context(), id)
	if err != nil {
		mapAndWriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

// ListMyBadges — GET /api/me/badges (UC15).
func (h *ProfileHandler) ListMyBadges(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	obtenus, err := h.profile.ListBadges(r.Context(), user.ID)
	if err != nil {
		mapAndWriteError(w, err)
		return
	}
	tous, err := h.profile.AllBadges(r.Context())
	if err != nil {
		mapAndWriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"obtenus": obtenus, "tous": tous})
}

// Leaderboard — GET /api/classement (UC17).
func (h *ProfileHandler) Leaderboard(w http.ResponseWriter, r *http.Request) {
	list, err := h.profile.Leaderboard(r.Context())
	if err != nil {
		mapAndWriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// LeaderboardForGroup — GET /api/groupes/{id}/classement (UC20).
func (h *ProfileHandler) LeaderboardForGroup(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	groupeID, err := idFromURL(r, "id")
	if err != nil {
		mapAndWriteError(w, services.ErrValidation)
		return
	}
	list, err := h.profile.LeaderboardForGroup(r.Context(), user.ID, groupeID)
	if err != nil {
		mapAndWriteError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}
