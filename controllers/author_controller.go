package controllers

import (
	"context"
	"digitalLibrary/config"
	"digitalLibrary/models"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateAuthor(firstName, familyName, dateOfBirthStr, dateOfDeathStr string) (interface{}, error) {
	if isZeroValue(firstName) == true || isZeroValue(familyName) == true {
		return false, errors.New("first_name, family_name or date_of_birth missing")
	}

	// dateOfBirth, err := time.Parse(time.RFC3339, dateOfBirthStr)

	// if err != nil {
	// 	return false, err
	// }

	author := models.Author{ID: primitive.NewObjectID(), FirstName: firstName, FamilyName: familyName}
	var err error
	var dateOfBirth time.Time
	if isZeroValue(dateOfBirthStr) == false {
		dateOfBirth, err = time.Parse(time.RFC3339, dateOfBirthStr)
		if err != nil {
			return false, err
		}
		author.DateOfBirth = dateOfBirth
	}

	var dateOfDeath time.Time
	if isZeroValue(dateOfDeathStr) == false {
		dateOfDeath, err = time.Parse(time.RFC3339, dateOfDeathStr)
		if err != nil {
			return false, err
		}
		author.DateOfDeath = dateOfDeath
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	insertResult, err := config.DB.Collection("authors").InsertOne(ctx, author)

	if err != nil {
		return false, err
	}
	newCustomLog("info", fmt.Sprintf("CreateAuthor book instance created %s", insertResult), make([]byte, 0), make([]byte, 0), err)
	return insertResult.InsertedID, nil
}

func GetAuthor(id string) (models.Author, error) {
	var result models.Author
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return result, err
	}

	err = config.DB.Collection("authors").FindOne(ctx, bson.M{"_id": oid}).Decode(&result)

	if err != nil {
		return result, err
	}

	return result, nil
}

func GetAuthors(limit, page int64) ([]models.Author, error) {
	var result []models.Author

	if (limit > 0) == false {
		limit = 5
	}
	if (page > 0) == false {
		page = 1
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	options := options.Find()
	options.SetSort(bson.D{{"first_name", 1}})
	options.SetSkip((page - 1) * limit)
	options.SetLimit(limit)

	cursor, err := config.DB.Collection("authors").Find(ctx, bson.M{}, options)
	if err != nil {
		return result, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return result, err
		//panic(err)
	}

	return result, nil
}

func GetTotalAuthorsCount(query bson.M) (int64, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	count, err := config.DB.Collection("authors").CountDocuments(ctx, query, nil)
	return count, err
}

func DeleteAuthor(id string) (bool, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return false, err
	}

	books, err := GetBooksByAuthorId(id)

	if err != nil {
		return false, err
	}

	if len(books) > 0 {
		return false, errors.New("please delete books before deleting author")
	}

	_, err = config.DB.Collection("authors").DeleteOne(ctx, bson.M{"_id": oid})

	if err != nil {
		return false, err
	}

	return true, nil
}

func UpdateAuthor(id string, author map[string]interface{}) (models.Author, error) {
	var authorDoc models.Author
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return authorDoc, err
	}
	query := bson.M{}
	query["_id"] = oid

	updateQuery := bson.M{}
	if author["first_name"] != nil {
		updateQuery["first_name"] = author["first_name"].(string)
	}
	if author["family_name"] != nil {
		updateQuery["family_name"] = author["family_name"].(string)
	}
	if author["date_of_birth"] != nil {
		dateOfBirth, err := time.Parse(time.RFC3339, author["date_of_birth"].(string))
		if err != nil {
			return authorDoc, err
		}
		updateQuery["date_of_birth"] = dateOfBirth
	}
	if author["date_of_death"] != nil {
		dateOfDeath, err := time.Parse(time.RFC3339, author["date_of_death"].(string))
		if err != nil {
			return authorDoc, err
		}
		updateQuery["date_of_death"] = dateOfDeath
	}

	if len(updateQuery) == 0 {
		return authorDoc, errors.New("nothing to update")
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	err = config.DB.Collection("authors").FindOneAndUpdate(ctx, query, bson.M{"$set": updateQuery}, &opt).Decode(&authorDoc)

	if err != nil {
		return authorDoc, err
	}

	return authorDoc, nil
}
