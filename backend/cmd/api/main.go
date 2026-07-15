// Command api démarre le serveur HTTP de "La Gueule de Bois".
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gueuledebois/backend/internal/config"
	"gueuledebois/backend/internal/handlers"
	"gueuledebois/backend/internal/repository/postgres"
	"gueuledebois/backend/internal/services"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("api: %v", err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(cfg.UploadDir, 0o755); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := postgres.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	// Repositories (accès Postgres via pgx).
	userRepo := postgres.NewUserRepo(pool)
	sessionRepo := postgres.NewSessionRepo(pool)
	soireeRepo := postgres.NewSoireeRepo(pool)
	photoRepo := postgres.NewPhotoRepo(pool)
	temoinInvitationRepo := postgres.NewTemoinInvitationRepo(pool)
	temoignageRepo := postgres.NewTemoignageRepo(pool)
	voteRepo := postgres.NewVoteRepo(pool)
	badgeRepo := postgres.NewBadgeRepo(pool)
	groupeRepo := postgres.NewGroupeRepo(pool)
	signalementRepo := postgres.NewSignalementRepo(pool)
	demandeAmiRepo := postgres.NewDemandeAmiRepo(pool)

	// Notifier : Resend si une clé API est fournie, sinon mock (dev local
	// sans dépendance externe, cf. spec).
	notifier := services.Notifier(services.NewMockNotifier())
	if cfg.ResendAPIKey != "" {
		notifier = services.NewResendNotifier(cfg.ResendAPIKey)
	} else {
		log.Println("api: RESEND_API_KEY absent -> notifications email simulées (MockNotifier)")
	}

	// Services (logique métier).
	authService := services.NewAuthService(userRepo, sessionRepo, cfg.SessionSecret)
	badgeService := services.NewBadgeService(badgeRepo, userRepo, notifier)
	scoringService := services.NewScoringService(userRepo, soireeRepo, temoignageRepo, voteRepo, badgeService)
	soireeService := services.NewSoireeService(soireeRepo, photoRepo, scoringService)
	temoignageService := services.NewTemoignageService(temoignageRepo, temoinInvitationRepo, soireeRepo, userRepo, scoringService, notifier)
	voteService := services.NewVoteService(voteRepo, temoignageRepo, soireeRepo, scoringService)
	profileService := services.NewProfileService(userRepo, badgeRepo, groupeRepo)
	signalementService := services.NewSignalementService(signalementRepo, temoignageRepo, soireeRepo, userRepo, scoringService, notifier)
	groupeService := services.NewGroupeService(groupeRepo)
	demandeAmiService := services.NewDemandeAmiService(demandeAmiRepo, userRepo, notifier)

	router := handlers.NewRouter(handlers.Deps{
		Auth:         authService,
		Soirees:      soireeService,
		Temoignages:  temoignageService,
		Votes:        voteService,
		Profile:      profileService,
		Signalements: signalementService,
		Groupes:      groupeService,
		Amis:         demandeAmiService,
		UploadDir:    cfg.UploadDir,
	})

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErrCh := make(chan error, 1)
	go func() {
		log.Printf("api: écoute sur :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrCh:
		return err
	case <-stop:
		log.Println("api: arrêt demandé, fermeture propre...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
