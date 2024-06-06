package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	Username        string     `json:"username" db:"username"`
	Email           string     `json:"email" db:"email"`
	ProfilePicture  string     `json:"profile_picture,omitempty" db:"profile_picture"`
	IsPremium       bool       `json:"is_premium" db:"is_premium"`
	EloPoints       int32      `json:"elo_points" db:"elo_points"`
	Password        string     `json:"password,omitempty" db:"password"`
	GoogleID        string     `json:"google_id,omitempty" db:"google_id"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty" db:"email_verified_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}
