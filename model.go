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
)

func init() {
	var ok bool
	dbConnectionString, ok = os.LookupEnv("DB_CONNECTION_STRING")
	if !ok {
		dbConnectionString = defautlDBConnectionString
	}
}

func OpenDB() (*gorm.DB, error) {
	return gorm.Open("mysql", dbConnectionString)
}

type UserRole struct {
	ID    int     `gorm:"primary_key; not null; auto_increment; column:id" json:"-"`
	Label *string `gorm:"not null; column:label" json:"label"`
}

type UserAddress struct {
	ID         int     `gorm:"primary_key; not null; auto_increment; column:id" json:"-"`
	Address    *string `gorm:"size:255;column:address" json:"address"`
	Province   *string `gorm:"size:255;column:province" json:"province"`
	City       *string `gorm:"size:255;column:city" json:"city"`
	Country    *string `gorm:"size:255;column:country" json:"country"`
	PostalCode *string `gorm:"size:255;column:postal_code" json:"postalCode"`
	UserId     int     `gorm:"not null; column:users_id" json:"-"`
}

type User struct {
	ID        int         `gorm:"primary_key; not null; auto_increment; column:id" json:"id"`
	Username  string      `gorm:"size: 255; not null; column:username" json:"username"`
	Email     string      `gorm:"size: 255; not null; column:email" json:"email"`
	CreatedAt *time.Time  `gorm:"not null; column:created_at" json:"createdAt"`
	UpdatedAt *time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	Role      UserRole    `gorm:"foreignkey:RoleId" json:"role"`
	RoleId    int         `gorm:"not null; column:user_roles_id" json:"-"`
	Address   UserAddress `gorm:"foreignkey:AddressId" json:"address"`
	AddressId int         `gorm:"not null; column:user_addresses_id" json:"-"`
}
