package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	Id                primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name              string             `json:"name" bson:"name"`
	Description       string             `json:"description" bson:"description"`
	ThumbnailImageUrl string             `json:"thumbnail" bson:"thumbnail"`
	Reviews           []Review           `json:"reviews" bson:"reviews"`
	RatingSum         float64            `json:"ratingSum" bson:"ratingSum"`
	RatingCount       int                `json:"ratingCount" bson:"ratingCount"`
}

func NewProduct(name, description, thumbnailImageUrl string) *Product {
	return &Product{
		Id:                primitive.NewObjectID(),
		Name:              name,
		Description:       description,
		ThumbnailImageUrl: thumbnailImageUrl,
		Reviews:           make([]Review, 0),
		RatingSum:         0,
		RatingCount:       0,
	}
}

func (p Product) Rating() float64 {
	if p.RatingCount == 0 {
		return 0
	}
	return p.RatingSum / float64(p.RatingCount)
}
