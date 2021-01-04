package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	ID      primitive.ObjectID   `json:"_id" bson:"_id"`
	Title   string               `json:"title" bson:"title"`
	Summary string               `json:"summary" bson:"summary"`
	Isbn    string               `json:"isbn" bson:"isbn"`
	Author  primitive.ObjectID   `json:"author" bson:"author"`
	Genre   []primitive.ObjectID `json:"genre" bson:"genre"`
}
