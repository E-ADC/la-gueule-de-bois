package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"gueuledebois/backend/internal/services"
)

// AuthHandler expose UC01 (inscription), UC02 (connexion), UC03
// (déconnexion) et un endpoint /me pour récupérer l'utilisateur courant.
type AuthHandler struct {
	auth *services.AuthService
}

func NewAuthHandler(auth *services.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

type registerRequest struct {
	Pseudo   string `json:"pseudo"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register — POST /api/auth/register (UC01).
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_input", "Corps de requête invalide.")
		return
	}

	req.Pseudo = strings.TrimSpace(req.Pseudo)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Pseudo == "" || req.Email == "" || len(req.Password) < 8 {
		writeError(w, http.StatusBadRequest, "invalid_input", "Pseudo, email et mot de passe (8 caractères min.) sont requis.")
		return
	}

	user, cookieValue, err := h.auth.Register(r.Context(), req.Pseudo, req.Email, req.Password)
	if err != nil {
		mapAndWriteError(w, err)
		return
	}

	setSessionCookie(w, cookieValue)
	writeJSON(w, http.StatusCreated, user)
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login — POST /api/auth/login (UC02).
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_input", "Corps de requête invalide.")
		return
	}

	user, cookieValue, err := h.auth.Login(r.Context(), strings.TrimSpace(strings.ToLower(req.Email)), req.Password)
	if err != nil {
		mapAndWriteError(w, err)
		return
	}

	setSessionCookie(w, cookieValue)
	writeJSON(w, http.StatusOK, user)
}

// Logout — POST /api/auth/logout (UC03).
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(services.SessionCookieName); err == nil {
		_ = h.auth.Logout(r.Context(), cookie.Value)
	}
	clearSessionCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

// Me — GET /api/auth/me : retourne l'utilisateur courant (route protégée).
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	writeJSON(w, http.StatusOK, user)
}

// setSessionCookie pose le cookie HttpOnly de session. Secure n'est pas
// activé ici pour simplifier le développement local en HTTP ; à activer
// (Secure: true) une fois le VPS derrière TLS (cf. README déploiement).
func setSessionCookie(w http.ResponseWriter, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     services.SessionCookieName,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(services.SessionTTL()),
	})
}

func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     services.SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
}
