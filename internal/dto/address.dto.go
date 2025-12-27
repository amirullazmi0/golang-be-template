package dto

// CreateAddressRequest represents create address request
type CreateAddressRequest struct {
	Label         string `json:"label" validate:"required,max=100"`
	RecipientName string `json:"recipient_name" validate:"required,max=255"`
	Phone         string `json:"phone" validate:"required,max=20"`
	Province      string `json:"province" validate:"required,max=100"`
	City          string `json:"city" validate:"required,max=100"`
	District      string `json:"district" validate:"required,max=100"`
	SubDistrict   string `json:"sub_district" validate:"required,max=100"`
	PostalCode    string `json:"postal_code" validate:"required,max=10"`
	FullAddress   string `json:"full_address" validate:"required"`
	IsPrimary     bool   `json:"is_primary"`
}

// UpdateAddressRequest represents update address request
type UpdateAddressRequest struct {
	ID            string `json:"id"`
	Label         string `json:"label" validate:"omitempty,max=100"`
	RecipientName string `json:"recipient_name" validate:"omitempty,max=255"`
	Phone         string `json:"phone" validate:"omitempty,max=20"`
	Province      string `json:"province" validate:"omitempty,max=100"`
	City          string `json:"city" validate:"omitempty,max=100"`
	District      string `json:"district" validate:"omitempty,max=100"`
	SubDistrict   string `json:"sub_district" validate:"omitempty,max=100"`
	PostalCode    string `json:"postal_code" validate:"omitempty,max=10"`
	FullAddress   string `json:"full_address" validate:"omitempty"`
	IsPrimary     *bool  `json:"is_primary"`
}

// AddressResponse represents address response
type AddressResponse struct {
	ID            string `json:"id"`
	UserID        string `json:"user_id"`
	Label         string `json:"label"`
	RecipientName string `json:"recipient_name"`
	Phone         string `json:"phone"`
	Province      string `json:"province"`
	City          string `json:"city"`
	District      string `json:"district"`
	SubDistrict   string `json:"sub_district"`
	PostalCode    string `json:"postal_code"`
	FullAddress   string `json:"full_address"`
	IsPrimary     bool   `json:"is_primary"`
	IsActive      bool   `json:"is_active"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}
