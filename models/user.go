package models

import (
	"errors"
	"fmt"
	"log"

	"github.com/MelchiSalins/go-auth/pkg/app"
	"github.com/jinzhu/gorm"
)

var (
	// ErrNotFound is returned when record is not found in the Datebase
	ErrNotFound = errors.New("models: Resource not found")
)

// User struct is GORM model to write to the database
// The values for this model are sourced from the response
// payload of the OAuth provider and are marshalled into this struct
type User struct {
	gorm.Model
	ISS           string `gorm:"type:varchar(10000)" json:"iss"`
	Email         string `gorm:"type:varchar(10000);unique;not null" json:"email"`
	EmailVerified bool   `gorm:"type:boolean" json:"email_verified"`
	Name          string `gorm:"type:varchar(10000)" json:"name"`
	Picture       string `gorm:"type:varchar(10000)" json:"picture"`
	GivenName     string `gorm:"type:varchar(10000)" json:"given_name"`
	FamilyName    string `gorm:"type:varchar(10000)" json:"family_name"`
	Locale        string `gorm:"type:varchar(10000)" json:"locale"`
	Iat           int    `gorm:"type:integer" json:"iat"`
	Exp           int    `gorm:"type:integer" jsong:"exp"`
}

// init runs migrations when modules is loaded.
func init() {
	// app.Init()
	// Start for DB migration to update schema
	// TODO: Migration should only run when needed, check if GORM can do this.

	fmt.Println("Running DB Migrations...")

	db, err := OpenDBConn(app.DBType, false)
	defer db.Close()
	if err != nil {
		log.Fatalln(err.Error())
	}
	db.AutoMigrate(&User{})
}

// NewUserService Returns a new UserService
func NewUserService() (*UserService, error) {
	db, err := OpenDBConn(app.DBType, false)
	if err != nil {
		return nil, err
	}

	return &UserService{
		db: db,
	}, nil
}

// UserService provides Methods to work with User Model.
type UserService struct {
	db *gorm.DB
}

// ByID Returns the first record of User of the matching ID
func (us *UserService) ByID(id int) (*User, error) {
	var user User

	err := us.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Close Closes the UserService database connection
func (us *UserService) Close() error {
	us.db.Close()
	return nil
}

// ExistOrCreate creates user from Tokenclaims if not already existing
// and returns a error
func (us *UserService) ExistOrCreate(u *User) error {
	if err := us.db.Create(&u).Error; err != nil {
		return err
	}

	return nil
}

// ByEmail Return a record based on Email ID
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User

	if err := us.db.Find(&user, &User{Email: email}).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
