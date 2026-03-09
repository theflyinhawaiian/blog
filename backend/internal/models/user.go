package models

import "time"

type User struct {
	ID          uint64    `db:"id" json:"id"`
	DisplayName string    `db:"display_name" json:"display_name"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type UserIdentity struct {
	ID             uint64  `db:"id" json:"id"`
	UserID         uint64  `db:"user_id" json:"user_id"`
	Provider       string  `db:"provider" json:"provider"`
	ProviderUserID string  `db:"provider_user_id" json:"provider_user_id"`
	Email          *string `db:"email" json:"email,omitempty"`
}
