package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type AccountModel struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	ClientID primitive.ObjectID `bson:"client_id"`
	Type     string             `bson:"type"`
	Balance  float64            `bson:"balance"`
}
