package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Review struct {
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Text     string             `json:"text" bson:"text"`
	Rating   int                `json:"rating" bson:"rating"`
	Reviewer User               `json:"user" bson:"user" mongo:"index"`
}

func NewReview(text string, rating int, reviewer User) *Review {
	return &Review{
		Id:       primitive.NewObjectID(),
		Text:     text,
		Rating:   rating,
		Reviewer: reviewer,
	}
}
