package usecase

import (
	"io"

	"github.com/amirullazmi0/kratify-backend/internal/dto"
	"github.com/amirullazmi0/kratify-backend/internal/repository"
)

type AttachmentUsecase interface {
	UploadImage(file io.Reader, fileName string) (*dto.AttachmentResponse, error)
	UploadDocument(file io.Reader, fileName string) (*dto.AttachmentResponse, error)
	UploadProductImage(file io.Reader, fileName string) (*dto.AttachmentResponse, error)
	UploadProfileImage(file io.Reader, fileName string, userID string) (*dto.AttachmentResponse, error)
}

type attachmentUsecase struct {
	repo repository.AttachmentRepository
}

func NewAttachmentUsecase(repo repository.AttachmentRepository) AttachmentUsecase {
	return &attachmentUsecase{
		repo: repo,
	}
}

func (u *attachmentUsecase) UploadImage(file io.Reader, fileName string) (*dto.AttachmentResponse, error) {
	return u.repo.UploadImage(file, fileName)
}

func (u *attachmentUsecase) UploadDocument(file io.Reader, fileName string) (*dto.AttachmentResponse, error) {
	return u.repo.UploadDocument(file, fileName)
}

func (u *attachmentUsecase) UploadProductImage(file io.Reader, fileName string) (*dto.AttachmentResponse, error) {
	return u.repo.UploadProductImage(file, fileName)
}

func (u *attachmentUsecase) UploadProfileImage(file io.Reader, fileName string, userID string) (*dto.AttachmentResponse, error) {
	return u.repo.UploadProfileImage(file, fileName, userID)
}
