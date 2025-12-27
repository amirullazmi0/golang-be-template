package usecase

import (
	"database/sql"
	"errors"
	"time"

	"github.com/amirullazmi0/kratify-backend/config"
	"github.com/amirullazmi0/kratify-backend/internal/dto"
	"github.com/amirullazmi0/kratify-backend/internal/middleware"
	"github.com/amirullazmi0/kratify-backend/internal/model"
	"github.com/amirullazmi0/kratify-backend/internal/repository"
)

type UserUsecase interface {
	Register(req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req *dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error)
	Logout(userID string) error
	GetProfile(userID string) (*dto.UserResponse, error)
	GetAllUsers() ([]dto.UserResponse, error)
	UpdateProfile(userID string, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	ChangePassword(userID string, req *dto.ChangePasswordRequest) error
	DeleteUser(userID string) error
}

type userUsecase struct {
	userRepo repository.UserRepository
	jwtCfg   *config.JWTConfig
}

// NewUserUsecase creates a new user usecase
func NewUserUsecase(userRepo repository.UserRepository, jwtCfg *config.JWTConfig) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
		jwtCfg:   jwtCfg,
	}
}

func (u *userUsecase) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if user already exists
	existingUser, _ := u.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Create new user
	user := &model.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	// Hash password
	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	// Save to database
	userID, err := u.userRepo.Create(user)
	if err != nil {
		return nil, err
	}
	user.ID = userID

	// Generate tokens
	accessToken, err := middleware.GenerateToken(user.ID, user.Email, u.jwtCfg)
	if err != nil {
		return nil, err
	}

	refreshToken, err := middleware.GenerateRefreshToken(user.ID, user.Email, u.jwtCfg)
	if err != nil {
		return nil, err
	}

	// Save refresh token to database
	refreshTokenExpiry := time.Now().Add(7 * 24 * time.Hour)
	if err := u.userRepo.SaveRefreshToken(user.ID, refreshToken, refreshTokenExpiry); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(u.jwtCfg.ExpiredHour * 3600),
		User: dto.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}

func (u *userUsecase) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Find user by email
	user, err := u.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Compare password
	if err := user.ComparePassword(req.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate tokens
	accessToken, err := middleware.GenerateToken(user.ID, user.Email, u.jwtCfg)
	if err != nil {
		return nil, err
	}

	refreshToken, err := middleware.GenerateRefreshToken(user.ID, user.Email, u.jwtCfg)
	if err != nil {
		return nil, err
	}

	// Save refresh token to database
	refreshTokenExpiry := time.Now().Add(7 * 24 * time.Hour)
	if err := u.userRepo.SaveRefreshToken(user.ID, refreshToken, refreshTokenExpiry); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(u.jwtCfg.ExpiredHour * 3600),
		User: dto.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}

func (u *userUsecase) GetProfile(userID string) (*dto.UserResponse, error) {
	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}, nil
}

func (u *userUsecase) GetAllUsers() ([]dto.UserResponse, error) {
	users, err := u.userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var response []dto.UserResponse
	for _, user := range users {
		response = append(response, dto.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		})
	}

	return response, nil
}

func (u *userUsecase) UpdateProfile(userID string, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}

	if err := u.userRepo.Update(user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}, nil
}

func (u *userUsecase) ChangePassword(userID string, req *dto.ChangePasswordRequest) error {
	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found")
		}
		return err
	}

	// Verify old password
	if err := user.ComparePassword(req.OldPassword); err != nil {
		return errors.New("invalid old password")
	}

	// Set new password
	user.Password = req.NewPassword
	if err := user.HashPassword(); err != nil {
		return err
	}

	return u.userRepo.Update(user)
}

func (u *userUsecase) DeleteUser(userID string) error {
	_, err := u.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found")
		}
		return err
	}

	return u.userRepo.Delete(userID)
}

func (u *userUsecase) RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	// Find user by refresh token
	user, err := u.userRepo.FindByRefreshToken(req.RefreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid or expired refresh token")
		}
		return nil, err
	}

	// Generate new access token
	accessToken, err := middleware.GenerateToken(user.ID, user.Email, u.jwtCfg)
	if err != nil {
		return nil, err
	}

	// Generate new refresh token
	newRefreshToken, err := middleware.GenerateRefreshToken(user.ID, user.Email, u.jwtCfg)
	if err != nil {
		return nil, err
	}

	// Update refresh token in database
	refreshTokenExpiry := time.Now().Add(7 * 24 * time.Hour)
	if err := u.userRepo.SaveRefreshToken(user.ID, newRefreshToken, refreshTokenExpiry); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(u.jwtCfg.ExpiredHour * 3600),
		User: dto.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}

func (u *userUsecase) Logout(userID string) error {
	// Verify user exists
	_, err := u.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found")
		}
		return err
	}

	// Clear refresh token
	return u.userRepo.ClearRefreshToken(userID)
}
