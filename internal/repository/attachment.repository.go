package repository

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"time"

	"github.com/amirullazmi0/kratify-backend/internal/dto"
	"github.com/amirullazmi0/kratify-backend/pkg/database"
	"github.com/amirullazmi0/kratify-backend/pkg/imagekit"
)

type AttachmentRepository interface {
	UploadImage(file io.Reader, fileName string) (*dto.AttachmentResponse, error)
	UploadDocument(file io.Reader, fileName string) (*dto.AttachmentResponse, error)
	UploadProductImage(file io.Reader, fileName string) (*dto.AttachmentResponse, error)
	UploadProfileImage(file io.Reader, fileName string, userID string) (*dto.AttachmentResponse, error)
}

type attachmentRepository struct {
	db       *sql.DB
	imageKit imagekit.ImageKitService
}

func NewAttachmentRepository(db *sql.DB, imageKit imagekit.ImageKitService) AttachmentRepository {
	return &attachmentRepository{
		db:       db,
		imageKit: imageKit,
	}
}

func (r *attachmentRepository) UploadImage(file io.Reader, fileName string) (*dto.AttachmentResponse, error) {
	resp, err := r.imageKit.UploadFile(context.Background(), file, fileName, "images")
	if err != nil {
		return nil, err
	}

	return &dto.AttachmentResponse{
		URL: resp.URL,
	}, nil
}

func (r *attachmentRepository) UploadDocument(file io.Reader, fileName string) (*dto.AttachmentResponse, error) {
	resp, err := r.imageKit.UploadFile(context.Background(), file, fileName, "documents")
	if err != nil {
		return nil, err
	}

	return &dto.AttachmentResponse{
		URL: resp.URL,
	}, nil
}

func (r *attachmentRepository) UploadProductImage(file io.Reader, fileName string) (*dto.AttachmentResponse, error) {
	resp, err := r.imageKit.UploadFile(context.Background(), file, fileName, "products")
	if err != nil {
		return nil, err
	}

	return &dto.AttachmentResponse{
		URL: resp.URL,
	}, nil
}

func (r *attachmentRepository) UploadProfileImage(file io.Reader, fileName string, userID string) (*dto.AttachmentResponse, error) {
	folder := fmt.Sprintf("profiles/%s", userID)
	resp, err := r.imageKit.UploadFile(context.Background(), file, fileName, folder)

	if err != nil {
		return nil, err
	}

	// Update user avatar
	_, err = database.NewUpdateBuilder("users").
		Set("avatar", resp.URL).
		Set("updated_at", time.Now()).
		Where("id = $1", userID).
		Execute(r.db)

	if err != nil {
		return nil, err
	}

	return &dto.AttachmentResponse{
		URL: resp.URL,
	}, nil
}
