package domain

type UserData struct {
	Username string `gorm:"type:nvarchar(200);not null;unique;" json:"username" validate:"required,min=6,max=200"`
	Name     string `gorm:"type:nvarchar(300);not null" json:"name" validate:"required,min=3,max=300"`
	Email    string `gorm:"type:nvarchar(300);not null;unique;" json:"email" validate:"required,email,max=300"`
	ShiftID  *uint  `gorm:"default:null" json:"-"`
	IsActive bool   `gorm:"default:true" json:"isActive"`
}
type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInResponse struct {
	UserData
	Token string `json:"token"`
}

type PasswordResetRequest struct {
	Username string `json:"username,omitempty"`
}

type NewUser struct {
	UserData
	Password string `json:"password" validate:"required,min=6"`
}

type ForgottenPasswordCode struct {
	Username string `json:"username" validate:"required,min=6,max=200"`
	OTP      string `json:"otp" validate:"required,min=6,max=6"`
}

type ResetForgottenPassword struct {
	Username        string `json:"username" validate:"required,min=6,max=200"`
	OTP             string `json:"otp" validate:"required,min=6,max=6"`
	Password        string `json:"password" validate:"required,min=6,max=200"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=6,max=200"`
}

type ChangePassword struct {
	CurrentPassword string `json:"currentPassword" validate:"required,min=6,max=200"`
	Password        string `json:"password" validate:"required,min=6,max=200"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=6,max=200"`
}

type UserTableCRUD struct { // use on web table GET /dashboard/settings/users
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	IsActive bool   `json:"isActive"`
	Profile  string `json:"profile"`
}

type CreateUserRequest struct {
	UserData
	ProfileID  uint   `json:"profileID" validate:"required"`
	Password   string `json:"password" validate:"required,min=6"`
	KitchenIDs []uint `json:"kitchenIDs"`
}

type EditableUser struct {
	ID uint `json:"id"`
	UserData
	ShiftID    *uint  `json:"shiftID"`
	ProfileID  uint   `json:"profileID" validate:"required"`
	KitchenIDs []uint `json:"kitchenIDs"`
}

type PermissionsByProfile struct {
	Writing bool `json:"writing"`
}

type ProfileWithDetails struct {
	ID          uint                   `json:"id"`
	Name        string                 `json:"name"`
	CreatedAt   string                 `json:"createdAt"`
	UpdatedAt   string                 `json:"updatedAt"`
	Permissions []PermissionsByProfile `json:"permissions"`
}

type CreateProfile struct {
	Name        string                 `json:"name" validate:"required,min=3,max=200"`
	Permissions []PermissionsByProfile `json:"permissions"`
}

type EditableProfile struct {
	ID          uint                   `json:"id"`
	Name        string                 `json:"name" validate:"required,min=3,max=200"`
	Permissions []PermissionsByProfile `json:"permissions"`
}
