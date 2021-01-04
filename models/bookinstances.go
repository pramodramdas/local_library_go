package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookInstance struct {
	ID      primitive.ObjectID `json:"_id" bson:"_id"`
	Book    primitive.ObjectID `json:"book" bson:"book"`
	Status  string             `json:"status" bson:"status"`
	Imprint string             `json:"imprint" bson:"imprint"`
	DueBack time.Time          `json:"due_back,omitempty" bson:"due_back,omitempty"`
}
