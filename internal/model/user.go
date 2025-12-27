package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string     `json:"id"`
	Email        string     `json:"email"`
	Password     string     `json:"-"`
	Name         string     `json:"name"`
	Role         string     `json:"role"`
	IsActive     bool       `json:"is_active"`
	RefreshToken *string    `json:"-"`
	TokenExpiry  *time.Time `json:"-"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	CreatedBy    *string    `json:"created_by,omitempty"`
	UpdatedBy    *string    `json:"updated_by,omitempty"`
	DeletedBy    *string    `json:"deleted_by,omitempty"`
}

// HashPassword hashes the user password
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// ComparePassword compares a password with the hashed password
func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
