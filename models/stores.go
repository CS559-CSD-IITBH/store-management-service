package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Store struct {
	StoreID     primitive.ObjectID `bson:"_id,omitempty"`
	MerchantID  uint               `bson:"merchant_id" json:"merchant_id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Available   bool               `bson:"available" json:"available"`
	Items       []Item             `bson:"items,omitempty" json:"items"`
}

type Item struct {
	ItemID      primitive.ObjectID `bson:"_id,omitempty"`
	StoreID     primitive.ObjectID `bson:"store_id" json:"store_id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Available   bool               `bson:"available" json:"available"`
	Price       float64            `bson:"price" json:"price"`
}
