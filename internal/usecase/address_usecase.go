package usecase

import (
	"github.com/amirullazmi0/kratify-backend/config"
	"github.com/amirullazmi0/kratify-backend/internal/dto"
	"github.com/amirullazmi0/kratify-backend/internal/repository"
)

type AddressUsecase interface {
	GetAddressByAuth(userID string) ([]dto.AddressResponse, error)
	GetAddressById(id string) (dto.AddressResponse, error)
	CreateAddress(userID string, body *dto.CreateAddressRequest) (dto.AddressResponse, error)
	UpdateAddress(userID string, body *dto.UpdateAddressRequest) (dto.AddressResponse, error)
	DeleteAddress(userID string, addressID string) error
}

type addressUsecase struct {
	addressRepo repository.AddressRepository
}

func NewAddressUsecase(addressRepo repository.AddressRepository, wtCfg *config.JWTConfig) AddressUsecase {
	return &addressUsecase{addressRepo: addressRepo}
}

func (u *addressUsecase) GetAddressByAuth(userID string) ([]dto.AddressResponse, error) {
	address := []dto.AddressResponse{}
	addresses, err := u.addressRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	for _, a := range addresses {
		address = append(address, dto.AddressResponse{
			ID:            a.ID,
			UserID:        a.UserID,
			Label:         a.Label,
			RecipientName: a.RecipientName,
			Phone:         a.Phone,
			Province:      a.Province,
			City:          a.City,
			District:      a.District,
			SubDistrict:   a.SubDistrict,
			PostalCode:    a.PostalCode,
			FullAddress:   a.FullAddress,
			IsPrimary:     a.IsPrimary,
			IsActive:      a.IsActive,
			CreatedAt:     a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     a.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return address, nil
}

func (u *addressUsecase) GetAddressById(id string) (dto.AddressResponse, error) {
	address, err := u.addressRepo.FindByID(id)
	if err != nil {
		return dto.AddressResponse{}, err
	}

	return dto.AddressResponse{
		ID:            address.ID,
		UserID:        address.UserID,
		Label:         address.Label,
		RecipientName: address.RecipientName,
		Phone:         address.Phone,
		Province:      address.Province,
		City:          address.City,
		District:      address.District,
		SubDistrict:   address.SubDistrict,
		PostalCode:    address.PostalCode,
		FullAddress:   address.FullAddress,
		IsPrimary:     address.IsPrimary,
		IsActive:      address.IsActive,
		CreatedAt:     address.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     address.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (u *addressUsecase) CreateAddress(userID string, body *dto.CreateAddressRequest) (dto.AddressResponse, error) {
	address, err := u.addressRepo.Create(userID, body)
	if err != nil {
		return dto.AddressResponse{}, err
	}

	return dto.AddressResponse{
		ID:            address.ID,
		UserID:        address.UserID,
		Label:         address.Label,
		RecipientName: address.RecipientName,
		Phone:         address.Phone,
		Province:      address.Province,
		City:          address.City,
		District:      address.District,
		SubDistrict:   address.SubDistrict,
		PostalCode:    address.PostalCode,
		FullAddress:   address.FullAddress,
		IsPrimary:     address.IsPrimary,
		IsActive:      address.IsActive,
		CreatedAt:     address.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     address.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (u *addressUsecase) UpdateAddress(userID string, body *dto.UpdateAddressRequest) (dto.AddressResponse, error) {
	address, err := u.addressRepo.Update(userID, body)
	if err != nil {
		return dto.AddressResponse{}, err
	}

	return dto.AddressResponse{
		ID:            address.ID,
		UserID:        address.UserID,
		Label:         address.Label,
		RecipientName: address.RecipientName,
		Phone:         address.Phone,
		Province:      address.Province,
		City:          address.City,
		District:      address.District,
		SubDistrict:   address.SubDistrict,
		PostalCode:    address.PostalCode,
		FullAddress:   address.FullAddress,
		IsPrimary:     address.IsPrimary,
		IsActive:      address.IsActive,
		CreatedAt:     address.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     address.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (u *addressUsecase) DeleteAddress(userID string, addressID string) error {
	// Verify address belongs to user
	address, err := u.addressRepo.FindByID(addressID)
	if err != nil {
		return err
	}

	// Check ownership
	if address.UserID != userID {
		return err
	}

	return u.addressRepo.Delete(addressID)
}
