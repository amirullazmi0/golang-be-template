package repository

import (
	"database/sql"
	"time"

	"github.com/amirullazmi0/kratify-backend/internal/dto"
	"github.com/amirullazmi0/kratify-backend/internal/model"
	"github.com/amirullazmi0/kratify-backend/pkg/database"
)

type AddressRepository interface {
	Create(userID string, address *dto.CreateAddressRequest) (model.Address, error)
	FindByID(id string) (*model.Address, error)
	FindByUserID(userID string) ([]model.Address, error)
	Update(userID string, address *dto.UpdateAddressRequest) (model.Address, error)
	Delete(id string) error
}

type addressRepository struct {
	db *sql.DB
}

func NewAddressRepository(db *sql.DB) AddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) Create(userID string, address *dto.CreateAddressRequest) (model.Address, error) {
	id, err := database.NewInsertBuilder("addresses").
		Set("user_id", userID).
		Set("label", address.Label).
		Set("recipient_name", address.RecipientName).
		Set("phone", address.Phone).
		Set("province", address.Province).
		Set("city", address.City).
		Set("district", address.District).
		Set("sub_district", address.SubDistrict).
		Set("postal_code", address.PostalCode).
		Set("full_address", address.FullAddress).
		Set("is_primary", address.IsPrimary).
		Execute(r.db)

	if err != nil {
		return model.Address{}, err
	}

	// Return created address
	return model.Address{
		ID:            id,
		UserID:        userID,
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
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

func (r *addressRepository) FindByID(id string) (*model.Address, error) {
	query, args := database.NewQueryBuilder("addresses").
		Where("id = $1", id).
		Where("deleted_at IS NULL").
		Limit(1).
		Build()

	var address model.Address
	err := r.db.QueryRow(query, args...).Scan(
		&address.ID,
		&address.UserID,
		&address.Label,
		&address.RecipientName,
		&address.Phone,
		&address.Province,
		&address.City,
		&address.District,
		&address.SubDistrict,
		&address.PostalCode,
		&address.FullAddress,
		&address.IsPrimary,
		&address.IsActive,
		&address.CreatedAt,
		&address.UpdatedAt,
		&address.DeletedAt,
		&address.CreatedBy,
		&address.UpdatedBy,
		&address.DeletedBy,
	)

	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, err
	}

	return &address, nil
}

func (r *addressRepository) FindByUserID(userID string) ([]model.Address, error) {
	query, args := database.NewQueryBuilder("addresses").
		Where("user_id = $1", userID).
		Where("deleted_at IS NULL").
		OrderBy("is_primary DESC, created_at DESC").
		Build()

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []model.Address
	for rows.Next() {
		var address model.Address
		err := rows.Scan(
			&address.ID,
			&address.UserID,
			&address.Label,
			&address.RecipientName,
			&address.Phone,
			&address.Province,
			&address.City,
			&address.District,
			&address.SubDistrict,
			&address.PostalCode,
			&address.FullAddress,
			&address.IsPrimary,
			&address.IsActive,
			&address.CreatedAt,
			&address.UpdatedAt,
			&address.DeletedAt,
			&address.CreatedBy,
			&address.UpdatedBy,
			&address.DeletedBy,
		)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}

func (r *addressRepository) Update(userID string, address *dto.UpdateAddressRequest) (model.Address, error) {
	// Build update dynamically based on non-nil fields
	builder := database.NewUpdateBuilder("addresses")

	if address.Label != "" {
		builder.Set("label", address.Label)
	}
	if address.RecipientName != "" {
		builder.Set("recipient_name", address.RecipientName)
	}
	if address.Phone != "" {
		builder.Set("phone", address.Phone)
	}
	if address.Province != "" {
		builder.Set("province", address.Province)
	}
	if address.City != "" {
		builder.Set("city", address.City)
	}
	if address.District != "" {
		builder.Set("district", address.District)
	}
	if address.SubDistrict != "" {
		builder.Set("sub_district", address.SubDistrict)
	}
	if address.PostalCode != "" {
		builder.Set("postal_code", address.PostalCode)
	}
	if address.FullAddress != "" {
		builder.Set("full_address", address.FullAddress)
	}
	if address.IsPrimary != nil {
		builder.Set("is_primary", *address.IsPrimary)
	}

	builder.Set("updated_at", time.Now())

	_, err := builder.
		Where("id = $1", address.ID).
		Where("user_id = $2", userID).
		Execute(r.db)
	if err != nil {
		return model.Address{}, err
	}

	// Fetch updated address
	updatedAddress, err := r.FindByID(address.ID)
	if err != nil {
		return model.Address{}, err
	}

	return *updatedAddress, nil
}

func (r *addressRepository) Delete(id string) error {
	// Soft delete
	_, err := database.NewUpdateBuilder("addresses").
		Set("deleted_at", time.Now()).
		Set("updated_at", time.Now()).
		Where("id = $1", id).
		Execute(r.db)

	return err
}
