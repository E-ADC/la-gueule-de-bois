package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"gueuledebois/backend/internal/services"
)

// Deps regroupe toutes les dépendances nécessaires au routeur : services
// métier déjà câblés à leurs repositories (cf. cmd/api/main.go).
type Deps struct {
	Auth         *services.AuthService
	Soirees      *services.SoireeService
	Temoignages  *services.TemoignageService
	Votes        *services.VoteService
	Profile      *services.ProfileService
	Signalements *services.SignalementService
	UploadDir    string
}

// NewRouter construit le routeur chi complet de l'API.
func NewRouter(deps Deps) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	authH := NewAuthHandler(deps.Auth)
	soireeH := NewSoireeHandler(deps.Soirees, deps.UploadDir)
	temoignageH := NewTemoignageHandler(deps.Temoignages)
	voteH := NewVoteHandler(deps.Votes)
	profileH := NewProfileHandler(deps.Profile)
	signalementH := NewSignalementHandler(deps.Signalements, deps.Temoignages)

	requireAuth := RequireAuth(deps.Auth)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/api", func(r chi.Router) {
		// UC01/02/03 — auth par sessions cookie HttpOnly (pas de JWT).
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authH.Register)
			r.Post("/login", authH.Login)
			r.Post("/logout", authH.Logout)
			r.With(requireAuth).Get("/me", authH.Me)
		})

		// UC05 — profil public (lecture, pas d'auth stricte requise mais on
		// la garde simple : la spec ne détaille pas de mode "visiteur").
		r.With(requireAuth).Get("/users/{id}", profileH.GetPublicProfile)

		// UC17 — classement global.
		r.With(requireAuth).Get("/classement", profileH.Leaderboard)
		// UC20 — classement restreint à un groupe.
		r.With(requireAuth).Get("/groupes/{id}/classement", profileH.LeaderboardForGroup)

		// UC15 — mes badges.
		r.With(requireAuth).Get("/me/badges", profileH.ListMyBadges)

		// UC06/07/08/10 — CRUD soirées + historique + photos.
		r.Route("/soirees", func(r chi.Router) {
			r.Use(requireAuth)
			r.Get("/", soireeH.ListMine)
			r.Post("/", soireeH.Create)
			r.Get("/{id}", soireeH.Get)
			r.Put("/{id}", soireeH.Update)
			r.Delete("/{id}", soireeH.Delete)
			r.Post("/{id}/photos", soireeH.UploadPhoto)

			// UC09 — inviter un témoin.
			r.Post("/{id}/temoins", temoignageH.InviteTemoin)
			// UC11 — ajouter/lister les témoignages d'une soirée.
			r.Post("/{id}/temoignages", temoignageH.Add)
			r.Get("/{id}/temoignages", temoignageH.ListBySoiree)
		})

		// UC12 — voter (swipe) sur un témoignage.
		r.With(requireAuth).Post("/temoignages/{id}/votes", voteH.Cast)

		// UC13 — signaler un témoignage.
		r.With(requireAuth).Post("/temoignages/{id}/signalements", signalementH.Report)

		// UC22 — traiter un signalement (acteur Modérateur uniquement).
		r.Route("/signalements", func(r chi.Router) {
			r.Use(requireAuth, RequireModerator)
			r.Get("/", signalementH.ListEnAttente)
			r.Post("/{id}/traiter", signalementH.Traiter)
		})

		registerTODORoutes(r)
	})

	// Photos servies en statique en local (en prod, nginx s'en charge — cf.
	// spec déploiement). Pratique pour tester le squelette sans nginx.
	fileServer := http.FileServer(http.Dir(deps.UploadDir))
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", fileServer))

	return r
}

// registerTODORoutes documente les cas d'utilisation non couverts par ce
// squelette (UC18, UC19, UC21) : modèles, migrations et repositories sont
// prêts (cf. internal/repository, internal/models,
// migrations/000001_init_schema.up.sql), il reste à écrire
// service+handler sur le même modèle que les routes ci-dessus.
func registerTODORoutes(r chi.Router) {
	todo := func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotImplemented, "not_implemented", "Non implémenté dans ce squelette (TODO).")
	}

	// TODO UC18 — Créer un groupe.
	r.Post("/groupes", todo)
	// TODO UC19 — Rejoindre un groupe.
	r.Post("/groupes/{id}/membres", todo)

	// TODO UC21 — Envoyer une demande d'ami.
	r.Post("/amis/demandes", todo)
}
