package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ClientModel struct {
	Id    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name"`
	Email string             `bson:"email"`
}
