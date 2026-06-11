package user

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	FullName        string     `json:"full_name"`
	Email           string     `json:"email"`
	Phone           string     `json:"phone"`
	Password        string     `json:"-"`
	AvatarURL       string     `json:"avatar_url"`
	Status          string     `json:"status"`
	Role            string     `json:"role"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	PhoneVerifiedAt *time.Time `json:"phone_verified_at"`
	LastLoginAt     *time.Time `json:"last_login_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
}
