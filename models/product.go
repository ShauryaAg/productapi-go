package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	Id                primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name              string             `json:"name" bson:"name"`
	Description       string             `json:"description" bson:"description"`
	ThumbnailImageUrl string             `json:"thumbnail" bson:"thumbnail"`
	Reviews           []Review           `json:"reviews" bson:"reviews"`
}
