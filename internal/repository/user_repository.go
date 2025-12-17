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
	FindAll() ([]model.User, error)
	Update(user *model.User) error
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
		Set("is_active", true).
		Execute(r.db)
	
	return id, err
}

func (r *userRepository) FindByID(id string) (*model.User, error) {
	query, args := database.NewQueryBuilder("users").
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
		&user.RefreshToken,
		&user.TokenExpiry,
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
		&user.RefreshToken,
		&user.TokenExpiry,
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
			&user.RefreshToken,
			&user.TokenExpiry,
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
