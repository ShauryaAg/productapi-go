package models

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Password string             `json:"password" bson:"password"`
	Email    string             `json:"email" bson:"email" mongo:"index,unique"`
	IsActive bool               `json:"is_active" bson:"is_active"`
}

func NewUser(name, email, password string) *User {
	user := &User{
		Id:       primitive.NewObjectID(),
		Name:     name,
		Email:    email,
		Password: password,
	}

	user.HashPassword()

	return user
}

func (u *User) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return err
	}

	u.Password = string(hash)
	return nil
}

func (u *User) VerifyPassword(attempt string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(attempt))
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func (u *User) SetPassword(new string) {
	u.Password = new
	u.HashPassword()
}
