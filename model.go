package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	defaultMySQLUser     = "iversoft_user"
	defaultMySQLPassword = "password"
	defaultMySQLHost     = ""
	defaultMySQLDBName   = "iversoft"
)

var (
	dbConnectionString string
	isDebug            bool
)

func init() {
	var (
		ok            bool
		mysqlUser     string
		mysqlPassword string
		mysqlHost     string
		mysqlDBName   string
	)

	mysqlUser, ok = os.LookupEnv("MYSQL_USER")
	if !ok {
		mysqlUser = defaultMySQLUser
	}

	mysqlPassword, ok = os.LookupEnv("MYSQL_PASSWORD")
	if !ok {
		mysqlPassword = defaultMySQLPassword
	}

	mysqlHost, ok = os.LookupEnv("MYSQL_HOST")
	if !ok {
		mysqlHost = defaultMySQLHost
	}

	mysqlDBName, ok = os.LookupEnv("MYSQL_DB_NAME")
	if !ok {
		mysqlDBName = defaultMySQLDBName
	}

	dbConnectionString = fmt.Sprintf("%s:%s@%s/%s?parseTime=True", mysqlUser, mysqlPassword, mysqlHost, mysqlDBName)

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

// TableName overrides `gorm.TableName`` for UserRole
func (s *UserRole) TableName(db *gorm.DB) string {
	return "user_roles"
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

// TableName overrides `gorm.TableName`` for UserAddress
func (s *UserAddress) TableName(db *gorm.DB) string {
	return "user_addresses"
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

// TableName overrides `gorm.TableName`` for User
func (s *User) TableName(db *gorm.DB) string {
	return "users"
}
