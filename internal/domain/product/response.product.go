package product

import "time"

type ProductResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	CreatedBy   string     `json:"createdBy"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
	UpdatedBy   *string    `json:"updatedBy,omitempty"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
	DeletedBy   *string    `json:"deletedBy,omitempty"`
}

type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
}
