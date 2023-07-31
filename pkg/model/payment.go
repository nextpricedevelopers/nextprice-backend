package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type PaymentOptions struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	PaymentName string             `bson:"payment_name" json:"payment_name"`
	Enabled     bool               `bson:"enabled" json:"enabled"`
	CreatedAt   string             `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt   string             `bson:"updated_at" json:"updated_at,omitempty"`
}