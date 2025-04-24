package domain

import (
	"time"

	"gorm.io/gorm"
)

type Database string

const (
	DBBetaAutopartes Database = "beta_autopartes"
)

type WebSocketWarehouses struct {
	ID            int       `gorm:"primary_key;auto_increment" json:"id"`
	SucursalID    int       `gorm:"column:id_sucursal;primary_key;auto_increment" json:"sucursalID"`
	LastUpdatedAt time.Time `gorm:"column:last_updated_at" json:"lastUpdatedAt"`
}

// table users - usuarios
type User struct {
	gorm.Model
	UserData
	OTP               string    `gorm:"type:nvarchar(6)" json:"otp"` // One Time Password
	OTPExpirationDate time.Time `gorm:"column otp_expiration_date" json:"otpExpirationDate"`
	Password          string    `gorm:"type:nvarchar(200);not null" json:"-" validate:"required,min=6,max=200"`
}

type Dev struct {
	gorm.Model
	Tag string `gorm:"type:nvarchar(200);not null" json:"tag" validate:"required,min=3,max=200"`
	IP  string `gorm:"type:nvarchar(200);not null" json:"ip" validate:"required,min=3,max=200"`
}
