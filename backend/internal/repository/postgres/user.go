package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"gueuledebois/backend/internal/models"
	"gueuledebois/backend/internal/repository"
)

// uniqueViolation code Postgres 23505 (contrainte UNIQUE) mappé vers
// repository.ErrConflict.
const pgUniqueViolation = "23505"

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation
}

// UserRepo implémente repository.UserRepository.
type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Create(ctx context.Context, u *models.User) error {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO users (pseudo, email, password_hash, avatar, bio, score, role)
		VALUES ($1, $2, $3, $4, $5, 0, 'user')
		RETURNING id, created_at`,
		u.Pseudo, u.Email, u.PasswordHash, u.Avatar, u.Bio,
	).Scan(&u.ID, &u.CreatedAt)
	if isUniqueViolation(err) {
		return repository.ErrConflict
	}
	return err
}

func (r *UserRepo) scanUser(row pgx.Row) (*models.User, error) {
	var u models.User
	err := row.Scan(&u.ID, &u.Pseudo, &u.Email, &u.PasswordHash, &u.Avatar, &u.Bio, &u.Score, &u.Role, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, pseudo, email, password_hash, avatar, bio, score, role, created_at
		FROM users WHERE id = $1`, id)
	return r.scanUser(row)
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, pseudo, email, password_hash, avatar, bio, score, role, created_at
		FROM users WHERE email = $1`, email)
	return r.scanUser(row)
}

func (r *UserRepo) GetByPseudo(ctx context.Context, pseudo string) (*models.User, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, pseudo, email, password_hash, avatar, bio, score, role, created_at
		FROM users WHERE pseudo = $1`, pseudo)
	return r.scanUser(row)
}

func (r *UserRepo) Update(ctx context.Context, u *models.User) error {
	ct, err := r.pool.Exec(ctx, `
		UPDATE users SET pseudo = $1, avatar = $2, bio = $3 WHERE id = $4`,
		u.Pseudo, u.Avatar, u.Bio, u.ID)
	if isUniqueViolation(err) {
		return repository.ErrConflict
	}
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *UserRepo) UpdateScore(ctx context.Context, userID int64, score int) error {
	ct, err := r.pool.Exec(ctx, `UPDATE users SET score = $1 WHERE id = $2`, score, userID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *UserRepo) ListLeaderboard(ctx context.Context, limit int) ([]models.User, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, pseudo, email, password_hash, avatar, bio, score, role, created_at
		FROM users ORDER BY score DESC, id ASC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return collectUsers(rows)
}

func (r *UserRepo) ListLeaderboardForGroup(ctx context.Context, groupeID int64, limit int) ([]models.User, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT u.id, u.pseudo, u.email, u.password_hash, u.avatar, u.bio, u.score, u.role, u.created_at
		FROM users u
		JOIN groupe_membres gm ON gm.user_id = u.id
		WHERE gm.groupe_id = $1
		ORDER BY u.score DESC, u.id ASC LIMIT $2`, groupeID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return collectUsers(rows)
}

func collectUsers(rows pgx.Rows) ([]models.User, error) {
	var out []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Pseudo, &u.Email, &u.PasswordHash, &u.Avatar, &u.Bio, &u.Score, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}
