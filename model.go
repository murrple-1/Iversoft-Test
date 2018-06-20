package main

import (
	"os"
	"time"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	defautlDBConnectionString = "iversoft_user:password@/iversoft?parseTime=True"
)

var (
	dbConnectionString string
	isDebug            bool
)

func init() {
	var ok bool
	dbConnectionString, ok = os.LookupEnv("DB_CONNECTION_STRING")
	if !ok {
		dbConnectionString = defautlDBConnectionString
	}

	if _isDebug, ok := os.LookupEnv("DEBUG"); ok {
		isDebug = (_isDebug == "true")
	} else {
		isDebug = false
	}
}

func openDB() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", dbConnectionString)
	if err == nil && isDebug {
		db = db.Debug()
	}
	return db, err
}

// UserRole represents an ORM mapping for `user_roles`
type UserRole struct {
	ID    int    `gorm:"primary_key; not null; auto_increment; column:id" json:"-"`
	Label string `gorm:"not null; column:label" json:"label"`
}

// UserAddress represents an ORM mapping for `user_addresses`
type UserAddress struct {
	ID         int     `gorm:"primary_key; not null; auto_increment; column:id" json:"-"`
	Address    *string `gorm:"size:255;column:address" json:"address"`
	Province   *string `gorm:"size:255;column:province" json:"province"`
	City       *string `gorm:"size:255;column:city" json:"city"`
	Country    *string `gorm:"size:255;column:country" json:"country"`
	PostalCode *string `gorm:"size:255;column:postal_code" json:"postalCode"`
}

// User represents an ORM mapping for `users`
type User struct {
	ID        int         `gorm:"primary_key; not null; auto_increment; column:id" json:"id"`
	Username  string      `gorm:"size: 255; not null; unique; column:username" json:"username"`
	Email     string      `gorm:"size: 255; not null; unique; column:email" json:"email"`
	CreatedAt time.Time   `gorm:"not null; column:created_at" json:"createdAt"`
	UpdatedAt *time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	Role      UserRole    `gorm:"foreignkey:RoleID" json:"role"`
	RoleID    int         `gorm:"not null; column:user_roles_id" json:"-"`
	Address   UserAddress `gorm:"foreignkey:AddressID" json:"address"`
	AddressID int         `gorm:"not null; column:user_addresses_id" json:"-"`
}
