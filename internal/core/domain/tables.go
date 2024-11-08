package domain

import (
	"time"

	"gorm.io/gorm"
)

// table permissions - permisos
type Permission struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"type:nvarchar(300);not null" json:"name" validate:"required,min=3,max=300"`
	Path string `gorm:"type:nvarchar(300);not null" json:"path" validate:"required,min=3,max=300"`
}

// table user_profiles - perfiles
type UserProfiles struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:nvarchar(300);not null" json:"name" validate:"required,min=3,max=300"`
	Slug      string         `gorm:"type:nvarchar(200);not null;unique" json:"slug" validate:"required,min=3,max=200"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// table profile_has_permissions - perfiles_has_permisos
type ProfilesHasPermissions struct {
	ProfileID    uint         `gorm:"not null" json:"profileId"`
	Profile      UserProfiles `gorm:"foreignKey:ProfileID;references:ID" json:"profile"`
	PermissionID uint         `gorm:"not null" json:"permissionId"`
	Permission   Permission   `gorm:"foreignKey:PermissionID;references:ID" json:"permission"`
	Writing      bool         `gorm:"not null" json:"writing"`
}

// table

// table users - usuarios
type User struct {
	gorm.Model
	UserData
	ProfileID         uint         `gorm:"not null" json:"-"`
	Profile           UserProfiles `gorm:"foreignKey:ProfileID;references:ID" json:"profile"`
	Shift             Shift        `gorm:"foreignKey:ShiftID;references:ID" json:"shift"`
	OTP               string       `gorm:"type:nvarchar(6)" json:"otp"` // One Time Password
	OTPExpirationDate time.Time    `gorm:"column otp_expiration_date" json:"otpExpirationDate"`
	Password          string       `gorm:"type:nvarchar(200);not null" json:"-" validate:"required,min=6,max=200"`
}

// table kitchen - cocinas
type Kitchen struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"type:nvarchar(300);not null" json:"name" validate:"required,min=3,max=300"`
}

// table usuarios_has_kitchens - usuarios_has_cocinas
type UsersHasKitchens struct {
	UserID    uint    `gorm:"not null" json:"userId"`
	User      User    `gorm:"foreignKey:UserID;references:ID" json:"user"`
	KitchenID uint    `gorm:"not null" json:"kitchenId"`
	Kitchen   Kitchen `gorm:"foreignKey:KitchenID;references:ID" json:"kitchen"`
}

// table shifts - turnos
type Shift struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"type:nvarchar(300);not null" json:"name" validate:"required,min=3,max=300"`
}

type Document struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"type:nvarchar(300);not null" json:"name" validate:"required,min=3,max=300"`
	Table     string         `gorm:"type:nvarchar(300);not null" json:"table" validate:"required,min=3,max=300"`
}

type DetailDocument struct {
	ID          uint     `gorm:"primaryKey" json:"id"`
	DocumentID  uint     `gorm:"not null" json:"-"`
	Document    Document `gorm:"foreignKey:DocumentID;references:ID" json:"-"`
	Field       string   `gorm:"type:nvarchar(300);not null" json:"field" validate:"required,min=3,max=300"`
	TypeField   string   `gorm:"type:nvarchar(300);not null" json:"typeField" validate:"required,min=3,max=300"`
	DocumentKey string   `gorm:"type:nvarchar(300);not null" json:"documentKey" validate:"required,min=3,max=300"`
}

type DocumentReports struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Name            string         `gorm:"type:nvarchar(300);not null" json:"name" validate:"required,min=3,max=300"`
	StoredProcedure string         `gorm:"type:nvarchar(300);not null" json:"-" validate:"required,min=3,max=300"`
}

type Dev struct {
	gorm.Model
	Tag string `gorm:"type:nvarchar(200);not null" json:"tag" validate:"required,min=3,max=200"`
	IP  string `gorm:"type:nvarchar(200);not null" json:"ip" validate:"required,min=3,max=200"`
}

/**
* Stored procedure field is a source of data table for the report
 */
type Reports struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Name            string         `gorm:"type:nvarchar(100);not null" json:"name" validate:"required,min=3,max=100"`
	StoredProcedure string         `gorm:"type:nvarchar(200);not null" json:"storedProcedure" validate:"required,min=5,max=200"`
}

type ChartReports struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Name            string         `gorm:"type:nvarchar(100);not null" json:"name" validate:"required,min=3,max=100"`
	StoredProcedure string         `gorm:"type:nvarchar(200);not null" json:"storedProcedure" validate:"required,min=5,max=200"`
	ReportID        uint           `gorm:"not null" json:"reportId"`
	Report          Reports        `gorm:"foreignKey:ReportID;references:ID" json:"report"`
}
