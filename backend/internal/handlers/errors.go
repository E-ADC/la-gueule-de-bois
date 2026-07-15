// Package handlers implémente la couche HTTP : routage, décodage JSON,
// validation des entrées et mapping des erreurs vers les codes HTTP et le
// format uniforme demandé par la spec : { "error": "...", "code": "..." }.
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"gueuledebois/backend/internal/repository"
	"gueuledebois/backend/internal/services"
)

// apiError est le format JSON uniforme de toutes les réponses d'erreur.
type apiError struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, apiError{Error: message, Code: code})
}

// mapAndWriteError traduit les erreurs sentinelles des services/repository
// vers le code HTTP et le code d'erreur attendus par la spec :
//
//	données invalides -> 400 ; non-propriétaire/non-membre -> 403 ;
//	ressource inexistante -> 404 ; doublon -> 409.
func mapAndWriteError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrValidation):
		writeError(w, http.StatusBadRequest, "invalid_input", "Données invalides.")
	case errors.Is(err, services.ErrForbidden):
		writeError(w, http.StatusForbidden, "forbidden", "Action non autorisée.")
	case errors.Is(err, repository.ErrNotFound):
		writeError(w, http.StatusNotFound, "not_found", "Ressource introuvable.")
	case errors.Is(err, repository.ErrConflict):
		writeError(w, http.StatusConflict, "conflict", "Cette action a déjà été effectuée.")
	case errors.Is(err, services.ErrIdentifiantsInvalides):
		writeError(w, http.StatusUnauthorized, "invalid_credentials", "Identifiants invalides.")
	case errors.Is(err, services.ErrSessionInvalide):
		writeError(w, http.StatusUnauthorized, "unauthenticated", "Session invalide ou expirée.")
	default:
		writeError(w, http.StatusInternalServerError, "internal_error", "Erreur interne.")
	}
}
