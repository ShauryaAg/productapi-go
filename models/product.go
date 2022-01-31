package models

import (
	validator "github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id                primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name              string             `json:"name" bson:"name" validate:"required"`
	Description       string             `json:"description" bson:"description" validate:"required"`
	ThumbnailImageUrl string             `json:"thumbnail" bson:"thumbnail" validate:"url"`
	Reviews           []Review           `json:"reviews" bson:"reviews"`
	RatingSum         float64            `json:"ratingSum" bson:"ratingSum"`
	RatingCount       int                `json:"ratingCount" bson:"ratingCount"`
}

func NewProduct(name, description, thumbnailImageUrl string) (*Product, error) {
	product := &Product{
		Id:                primitive.NewObjectID(),
		Name:              name,
		Description:       description,
		ThumbnailImageUrl: thumbnailImageUrl,
		Reviews:           make([]Review, 0),
		RatingSum:         0,
		RatingCount:       0,
	}

	v := validator.New()
	err := v.Struct(product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p Product) Rating() float64 {
	if p.RatingCount == 0 {
		return 0
	}
	return p.RatingSum / float64(p.RatingCount)
}
