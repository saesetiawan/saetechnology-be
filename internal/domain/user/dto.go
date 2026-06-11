package user

type CreateUserDto struct {
	FullName string `json:"full_name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=100"`
	Phone    string `json:"phone" validate:"required,min=8,max=20"`
	Role     string `json:"role" validate:"omitempty,oneof=admin customer"`
}

type UpdateProfileDto struct {
	FullName string `json:"full_name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,min=8,max=20"`
}

type ChangePasswordDto struct {
	CurrentPassword         string `json:"current_password" validate:"required"`
	NewPassword             string `json:"new_password" validate:"required,min=8,max=100"`
	NewPasswordConfirmation string `json:"new_password_confirmation" validate:"required"`
}

type ListUsersQuery struct {
	Page      int
	Limit     int
	Search    string
	SearchBy  string
	Role      string
	Status    string
	OrderBy   string
	OrderType string
}

type AdminCreateUserDto struct {
	FullName string `json:"full_name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,min=8,max=20"`
	Password string `json:"password" validate:"required,min=8,max=100"`
	Role     string `json:"role" validate:"required,oneof=admin customer"`
	Status   string `json:"status" validate:"required,oneof=active inactive"`
}

type AdminUpdateUserStatusDto struct {
	Status string `json:"status" validate:"required,oneof=active inactive"`
}
