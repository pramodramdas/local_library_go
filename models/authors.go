package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Author struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	DateOfBirth time.Time          `json:"date_of_birth,omitempty" bson:"date_of_birth,omitempty"`
	DateOfDeath time.Time          `json:"date_of_death,omitempty" bson:"date_of_death,omitempty"`
	FirstName   string             `json:"first_name" bson:"first_name"`
	FamilyName  string             `json:"family_name" bson:"family_name"`
}
