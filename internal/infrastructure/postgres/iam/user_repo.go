package iam

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/masterfabric/masterfabric_go_basic/internal/domain/iam/model"
	domainErr "github.com/masterfabric/masterfabric_go_basic/internal/shared/errors"
)

// UserRepo is the PostgreSQL implementation of domain UserRepository.
type UserRepo struct {
	db *pgxpool.Pool
}

// NewUserRepo creates a new UserRepo.
func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

const (
	sqlCreateUser = `
		INSERT INTO users (id, email, password_hash, display_name, avatar_url, bio, status, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	sqlFindUserByID = `
		SELECT id, email, password_hash, display_name, avatar_url, bio, status, role, created_at, updated_at
		FROM users WHERE id = $1`

	sqlFindUserByEmail = `
		SELECT id, email, password_hash, display_name, avatar_url, bio, status, role, created_at, updated_at
		FROM users WHERE email = $1`

	sqlUpdateUser = `
		UPDATE users
		SET email = $2, password_hash = $3, display_name = $4, avatar_url = $5, bio = $6,
		    status = $7, role = $8, updated_at = $9
		WHERE id = $1`

	sqlDeleteUser = `DELETE FROM users WHERE id = $1`

	sqlListUsers = `
		SELECT id, email, password_hash, display_name, avatar_url, bio, status, role, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	sqlCountUsers = `SELECT COUNT(*) FROM users`
)

// Create persists a new user record.
func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	_, err := r.db.Exec(ctx, sqlCreateUser,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.DisplayName,
		user.AvatarURL,
		user.Bio,
		string(user.Status),
		string(user.Role),
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("userRepo.Create: %w", err)
	}
	return nil
}

// FindByID retrieves a user by primary key.
func (r *UserRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	row := r.db.QueryRow(ctx, sqlFindUserByID, id)
	user, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.ErrUserNotFound
		}
		return nil, fmt.Errorf("userRepo.FindByID: %w", err)
	}
	return user, nil
}

// FindByEmail retrieves a user by email address.
func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	row := r.db.QueryRow(ctx, sqlFindUserByEmail, email)
	user, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.ErrUserNotFound
		}
		return nil, fmt.Errorf("userRepo.FindByEmail: %w", err)
	}
	return user, nil
}

// Update persists changes to an existing user.
func (r *UserRepo) Update(ctx context.Context, user *model.User) error {
	user.UpdatedAt = time.Now().UTC()
	tag, err := r.db.Exec(ctx, sqlUpdateUser,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.DisplayName,
		user.AvatarURL,
		user.Bio,
		string(user.Status),
		string(user.Role),
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("userRepo.Update: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domainErr.ErrUserNotFound
	}
	return nil
}

// Delete removes a user by ID.
func (r *UserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.db.Exec(ctx, sqlDeleteUser, id)
	if err != nil {
		return fmt.Errorf("userRepo.Delete: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domainErr.ErrUserNotFound
	}
	return nil
}

// List returns a paginated list of all users ordered by creation date descending.
func (r *UserRepo) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
	rows, err := r.db.Query(ctx, sqlListUsers, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("userRepo.List: %w", err)
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("userRepo.List scan: %w", err)
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

// CountAll returns the total number of users.
func (r *UserRepo) CountAll(ctx context.Context) (int, error) {
	var count int
	if err := r.db.QueryRow(ctx, sqlCountUsers).Scan(&count); err != nil {
		return 0, fmt.Errorf("userRepo.CountAll: %w", err)
	}
	return count, nil
}

// scanUser reads a single pgx row into a User entity.
func scanUser(row pgx.Row) (*model.User, error) {
	var (
		u      model.User
		status string
		role   string
	)
	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.DisplayName,
		&u.AvatarURL,
		&u.Bio,
		&status,
		&role,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	u.Status = model.UserStatus(status)
	u.Role = model.UserRole(role)
	return &u, nil
}
