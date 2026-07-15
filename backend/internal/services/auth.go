package services

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"gueuledebois/backend/internal/models"
	"gueuledebois/backend/internal/repository"
)

// ErrIdentifiantsInvalides est renvoyée par Login en cas de mauvais email
// ou mot de passe (UC02, exception "identifiants invalides -> accès refusé").
var ErrIdentifiantsInvalides = errors.New("auth: identifiants invalides")

// ErrSessionInvalide est renvoyée quand le cookie de session est absent,
// mal signé, inconnu en base, ou expiré.
var ErrSessionInvalide = errors.New("auth: session invalide")

// sessionTTL est la durée de vie d'une session (7 jours, valeur simple et
// raisonnable pour un projet école — non précisée par les fiches).
const sessionTTL = 7 * 24 * time.Hour

// AuthService gère l'inscription, la connexion et les sessions cookie
// HttpOnly (UC01/02/03), sans JWT (spec).
type AuthService struct {
	users    repository.UserRepository
	sessions repository.SessionRepository
	secret   []byte
}

func NewAuthService(users repository.UserRepository, sessions repository.SessionRepository, secret string) *AuthService {
	return &AuthService{users: users, sessions: sessions, secret: []byte(secret)}
}

// Register implémente UC01 : crée un compte si le pseudo/email est libre,
// hash le mot de passe (bcrypt) puis ouvre une session (connexion auto).
func (s *AuthService) Register(ctx context.Context, pseudo, email, password string) (*models.User, string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("auth: hash mot de passe: %w", err)
	}

	// Role fixé côté application (et non laissé au zero-value Go) pour que
	// l'utilisateur retourné immédiatement après inscription reflète bien
	// la valeur par défaut posée en base ('user', cf. migration 000001).
	u := &models.User{Pseudo: pseudo, Email: email, PasswordHash: string(hash), Role: "user"}
	if err := s.users.Create(ctx, u); err != nil {
		return nil, "", err // repository.ErrConflict si email/pseudo pris
	}

	cookie, err := s.openSession(ctx, u.ID)
	if err != nil {
		return nil, "", err
	}
	return u, cookie, nil
}

// Login implémente UC02 : vérifie les identifiants puis ouvre une session.
func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, string, error) {
	u, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, "", ErrIdentifiantsInvalides
		}
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, "", ErrIdentifiantsInvalides
	}

	cookie, err := s.openSession(ctx, u.ID)
	if err != nil {
		return nil, "", err
	}
	return u, cookie, nil
}

// Logout implémente UC03 : ferme la session côté serveur. cookieValue est
// la valeur brute lue dans le cookie (token + signature).
func (s *AuthService) Logout(ctx context.Context, cookieValue string) error {
	token, ok := s.verifyCookie(cookieValue)
	if !ok {
		return nil // rien à faire, cookie déjà invalide
	}
	return s.sessions.DeleteByToken(ctx, token)
}

// CurrentUser résout l'utilisateur associé à un cookie de session valide.
func (s *AuthService) CurrentUser(ctx context.Context, cookieValue string) (*models.User, error) {
	token, ok := s.verifyCookie(cookieValue)
	if !ok {
		return nil, ErrSessionInvalide
	}

	sess, err := s.sessions.GetByToken(ctx, token)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrSessionInvalide
		}
		return nil, err
	}
	if time.Now().After(sess.ExpiresAt) {
		_ = s.sessions.DeleteByToken(ctx, token)
		return nil, ErrSessionInvalide
	}

	return s.users.GetByID(ctx, sess.UserID)
}

// openSession crée une nouvelle session en base et retourne la valeur de
// cookie signée (token + HMAC(SESSION_SECRET, token)). La signature évite
// de taper la base pour des cookies grossièrement forgés/aléatoires.
func (s *AuthService) openSession(ctx context.Context, userID int64) (string, error) {
	token, err := randomToken()
	if err != nil {
		return "", fmt.Errorf("auth: génération token: %w", err)
	}

	sess := &models.Session{
		Token:     token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(sessionTTL),
	}
	if err := s.sessions.Create(ctx, sess); err != nil {
		return "", fmt.Errorf("auth: création session: %w", err)
	}

	return s.signCookie(token), nil
}

func (s *AuthService) signCookie(token string) string {
	return token + "." + s.sign(token)
}

func (s *AuthService) verifyCookie(cookieValue string) (token string, ok bool) {
	parts := strings.SplitN(cookieValue, ".", 2)
	if len(parts) != 2 {
		return "", false
	}
	token, sig := parts[0], parts[1]
	expected := s.sign(token)
	if !hmac.Equal([]byte(sig), []byte(expected)) {
		return "", false
	}
	return token, true
}

func (s *AuthService) sign(token string) string {
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(token))
	return hex.EncodeToString(mac.Sum(nil))
}

func randomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// SessionCookieName est le nom du cookie HttpOnly transportant la session.
const SessionCookieName = "gdb_session"

// SessionTTL expose la durée de vie de session pour la couche handlers
// (calcul de l'expiration du cookie HTTP).
func SessionTTL() time.Duration { return sessionTTL }
