package user

import "context"

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
	FindAll(ctx context.Context, query ListUsersQuery) ([]User, int64, error)
	Create(ctx context.Context, user *User) error
	Activate(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status string) (*User, error)
	UpdateProfile(ctx context.Context, id string, payload UpdateProfileDto) (*User, error)
	UpdatePassword(ctx context.Context, id string, passwordHash string) error
}
