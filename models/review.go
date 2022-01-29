package models

import (
	validator "github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" `
	Text     string             `json:"text" bson:"text" validate:"required"`
	Rating   int                `json:"rating" bson:"rating" validate:"required,min=1,max=5"`
	Reviewer User               `json:"user" bson:"user" mongo:"index" validate:"required"`
}

func NewReview(text string, rating int, reviewer User) (*Review, error) {

	review := &Review{
		Id:       primitive.NewObjectID(),
		Text:     text,
		Rating:   rating,
		Reviewer: reviewer,
	}

	v := validator.New()
	err := v.Struct(review)
	if err != nil {
		return nil, err
	}
	return review, nil
}
