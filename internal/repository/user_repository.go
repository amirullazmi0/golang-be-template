package repository

import (
	"database/sql"
	"time"

	"github.com/amirullazmi0/kratify-backend/internal/model"
	"github.com/amirullazmi0/kratify-backend/pkg/database"
)

type UserRepository interface {
	Create(user *model.User) (string, error)
	FindByID(id string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindByRefreshToken(refreshToken string) (*model.User, error)
	FindByVerificationToken(token string) (*model.User, error)
	FindAll() ([]model.User, error)
	Update(user *model.User) error
	SaveRefreshToken(userID string, refreshToken string, expiresAt time.Time) error
	ClearRefreshToken(userID string) error
	SaveVerificationToken(userID string, token string, expiresAt time.Time) error
	VerifyEmail(userID string) error
	Delete(id string) error
}

type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) (string, error) {
	id, err := database.NewInsertBuilder("users").
		Set("email", user.Email).
		Set("password", user.Password).
		Set("name", user.Name).
		Execute(r.db)

	return id, err
}

func (r *userRepository) FindByID(id string) (*model.User, error) {
	query, args := database.NewQueryBuilder("users").
		Select("id", "email", "password", "name", "role", "refresh_token", "token_expiry", "verification_token", "verification_expiry", "is_active", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "deleted_by").
		Where("id = $1", id).
		Where("deleted_at IS NULL").
		Limit(1).
		Build()

	var user model.User
	err := r.db.QueryRow(query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Role,
		&user.RefreshToken,
		&user.TokenExpiry,
		&user.VerificationToken,
		&user.VerificationExpiry,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
		&user.CreatedBy,
		&user.UpdatedBy,
		&user.DeletedBy,
	)

	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	query, args := database.NewQueryBuilder("users").
		Select("id", "email", "password", "name", "role", "refresh_token", "token_expiry", "verification_token", "verification_expiry", "is_active", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "deleted_by").
		Where("email = $1", email).
		Where("deleted_at IS NULL").
		Limit(1).
		Build()

	var user model.User
	err := r.db.QueryRow(query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Role,
		&user.RefreshToken,
		&user.TokenExpiry,
		&user.VerificationToken,
		&user.VerificationExpiry,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
		&user.CreatedBy,
		&user.UpdatedBy,
		&user.DeletedBy,
	)

	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindAll() ([]model.User, error) {
	query, args := database.NewQueryBuilder("users").
		Select("id", "email", "password", "name", "role", "refresh_token", "token_expiry", "verification_token", "verification_expiry", "is_active", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "deleted_by").
		Where("deleted_at IS NULL").
		OrderBy("created_at DESC").
		Build()

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.Name,
			&user.Role,
			&user.RefreshToken,
			&user.TokenExpiry,
			&user.VerificationToken,
			&user.VerificationExpiry,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.CreatedBy,
			&user.UpdatedBy,
			&user.DeletedBy,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) Update(user *model.User) error {
	_, err := database.NewUpdateBuilder("users").
		Set("name", user.Name).
		Set("password", user.Password).
		Set("updated_at", time.Now()).
		Where("id = $1", user.ID).
		Execute(r.db)

	return err
}

func (r *userRepository) Delete(id string) error {
	// Soft delete
	_, err := database.NewUpdateBuilder("users").
		Set("deleted_at", time.Now()).
		Set("updated_at", time.Now()).
		Where("id = $1", id).
		Execute(r.db)

	return err
}

func (r *userRepository) SaveRefreshToken(userID string, refreshToken string, expiresAt time.Time) error {
	_, err := database.NewUpdateBuilder("users").
		Set("refresh_token", refreshToken).
		Set("token_expiry", expiresAt).
		Set("updated_at", time.Now()).
		Where("id = $1", userID).
		Execute(r.db)

	return err
}

func (r *userRepository) ClearRefreshToken(userID string) error {
	_, err := database.NewUpdateBuilder("users").
		Set("refresh_token", nil).
		Set("token_expiry", nil).
		Set("updated_at", time.Now()).
		Where("id = $1", userID).
		Execute(r.db)

	return err
}

func (r *userRepository) SaveVerificationToken(userID string, token string, expiresAt time.Time) error {
	_, err := database.NewUpdateBuilder("users").
		Set("verification_token", token).
		Set("verification_expiry", expiresAt).
		Set("updated_at", time.Now()).
		Where("id = $1", userID).
		Execute(r.db)

	return err
}

func (r *userRepository) VerifyEmail(userID string) error {
	_, err := database.NewUpdateBuilder("users").
		Set("is_active", true).
		Set("verification_token", nil).
		Set("verification_expiry", nil).
		Set("updated_at", time.Now()).
		Where("id = $1", userID).
		Execute(r.db)

	return err
}

func (r *userRepository) FindByVerificationToken(token string) (*model.User, error) {
	query, args := database.NewQueryBuilder("users").
		Select("id", "email", "password", "name", "role", "refresh_token", "token_expiry", "verification_token", "verification_expiry", "is_active", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "deleted_by").
		Where("verification_token = $1", token).
		Where("deleted_at IS NULL").
		Where("verification_expiry > $2", time.Now()).
		Limit(1).
		Build()

	var user model.User
	err := r.db.QueryRow(query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Role,
		&user.RefreshToken,
		&user.TokenExpiry,
		&user.VerificationToken,
		&user.VerificationExpiry,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
		&user.CreatedBy,
		&user.UpdatedBy,
		&user.DeletedBy,
	)

	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByRefreshToken(refreshToken string) (*model.User, error) {
	query, args := database.NewQueryBuilder("users").
		Select("id", "email", "password", "name", "role", "refresh_token", "token_expiry", "verification_token", "verification_expiry", "is_active", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "deleted_by").
		Where("refresh_token = $1", refreshToken).
		Where("deleted_at IS NULL").
		Where("token_expiry > $2", time.Now()).
		Limit(1).
		Build()

	var user model.User
	err := r.db.QueryRow(query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Role,
		&user.RefreshToken,
		&user.TokenExpiry,
		&user.VerificationToken,
		&user.VerificationExpiry,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
		&user.CreatedBy,
		&user.UpdatedBy,
		&user.DeletedBy,
	)

	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
