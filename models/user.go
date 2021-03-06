package models

import (
	validator "github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name" validate:"required"`
	Password string             `json:"password" bson:"password" validate:"required"`
	Email    string             `json:"email" bson:"email" mongo:"unique" validate:"required,email"`
}

// NewUser creates a new user and hashes the password
// and returns an error if the user is invalid
func NewUser(name, email, password string) (*User, error) {
	user := &User{
		Id:       primitive.NewObjectID(),
		Name:     name,
		Email:    email,
		Password: password,
	}

	v := validator.New()
	err := v.Struct(user)
	if err != nil {
		return nil, err
	}

	user.HashPassword()
	return user, nil
}

// HashPassword hashes the user's password
func (u *User) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hash)
	return nil
}

// ComparePassword compares the user's password
// with the attempted password
func (u *User) VerifyPassword(attempt string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(attempt))
	if err != nil {
		return false
	}

	return true
}

// SetPassword sets the user's password
// and hashes it
func (u *User) SetPassword(new string) {
	u.Password = new
	u.HashPassword()
}
