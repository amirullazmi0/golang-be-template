package model

import "time"

type Address struct {
	ID            string     `json:"id"`
	UserID        string     `json:"user_id"`
	Label         string     `json:"label"`
	RecipientName string     `json:"recipient_name"`
	Phone         string     `json:"phone"`
	Province      string     `json:"province"`
	City          string     `json:"city"`
	District      string     `json:"district"`
	SubDistrict   string     `json:"sub_district"`
	PostalCode    string     `json:"postal_code"`
	FullAddress   string     `json:"full_address"`
	IsPrimary     bool       `json:"is_primary"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
	CreatedBy     *string    `json:"created_by,omitempty"`
	UpdatedBy     *string    `json:"updated_by,omitempty"`
	DeletedBy     *string    `json:"deleted_by,omitempty"`
}
