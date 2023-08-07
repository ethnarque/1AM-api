package data

import "time"

type TranslationKey struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
}
