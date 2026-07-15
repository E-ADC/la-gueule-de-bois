package handlers

import (
	"context"
	"net/http"

	"gueuledebois/backend/internal/models"
	"gueuledebois/backend/internal/services"
)

type ctxKey int

const userCtxKey ctxKey = iota

// RequireAuth résout l'utilisateur courant à partir du cookie de session
// et le place dans le contexte de la requête. Répond 401 si absent/invalide.
func RequireAuth(auth *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(services.SessionCookieName)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "unauthenticated", "Connexion requise.")
				return
			}

			user, err := auth.CurrentUser(r.Context(), cookie.Value)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "unauthenticated", "Session invalide ou expirée.")
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserFromContext récupère l'utilisateur placé par RequireAuth. Ne doit
// être appelé que sur des routes protégées par ce middleware.
func UserFromContext(ctx context.Context) *models.User {
	u, _ := ctx.Value(userCtxKey).(*models.User)
	return u
}

// RequireModerator restreint l'accès aux utilisateurs de rôle "moderator"
// (UC22, acteur Modérateur). Doit être chaîné après RequireAuth.
func RequireModerator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := UserFromContext(r.Context())
		if user == nil || user.Role != "moderator" {
			writeError(w, http.StatusForbidden, "forbidden", "Action réservée aux modérateurs.")
			return
		}
		next.ServeHTTP(w, r)
	})
}
