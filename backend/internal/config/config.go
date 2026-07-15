// Package config centralise la lecture de la configuration applicative
// à partir des variables d'environnement (12-factor).
package config

import (
	"fmt"
	"os"
)

// Config regroupe tous les paramètres nécessaires au démarrage de l'API.
type Config struct {
	Port          string
	DatabaseURL   string
	ResendAPIKey  string
	SessionSecret string
	UploadDir     string
}

// Load lit les variables d'environnement et applique des valeurs par défaut
// raisonnables pour le développement local. DATABASE_URL et SESSION_SECRET
// sont obligatoires : on préfère échouer au démarrage plutôt que de tourner
// dans un état incohérent.
func Load() (Config, error) {
	cfg := Config{
		Port:          getEnvDefault("PORT", "8080"),
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		ResendAPIKey:  os.Getenv("RESEND_API_KEY"),
		SessionSecret: os.Getenv("SESSION_SECRET"),
		UploadDir:     getEnvDefault("UPLOAD_DIR", "./uploads"),
	}

	if cfg.DatabaseURL == "" {
		return cfg, fmt.Errorf("config: DATABASE_URL est requis")
	}
	if cfg.SessionSecret == "" {
		return cfg, fmt.Errorf("config: SESSION_SECRET est requis")
	}

	return cfg, nil
}

func getEnvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
