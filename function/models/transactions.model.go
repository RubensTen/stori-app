package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionModel struct {
	Id              primitive.ObjectID `bson:"_id,omitempty"`
	AccountID       primitive.ObjectID `bson:"account_id"`
	Ammount         float64            `bson:"ammount"`
	TransactionDate string             `bson:"transaction_date"`
}
