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

func CreateGenre(genreName string) (interface{}, error) {
	if isZeroValue(genreName) == true {
		return false, errors.New("genreName missing")
	}

	genre := models.Genre{ID: primitive.NewObjectID(), Name: genreName}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	insertResult, err := config.DB.Collection("genres").InsertOne(ctx, genre)

	if err != nil {
		return false, err
	}
	newCustomLog("info", fmt.Sprintf("CreateAuthor book instance created %s", insertResult), make([]byte, 0), make([]byte, 0), err)
	return insertResult.InsertedID, nil
}

func GetGenre(id string) (models.Genre, error) {
	var result models.Genre
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return result, err
	}

	err = config.DB.Collection("genres").FindOne(ctx, bson.M{"_id": oid}).Decode(&result)

	if err != nil {
		return result, err
	}

	return result, nil
}

func GetGenres(limit, page int64) ([]models.Genre, error) {
	var result []models.Genre

	if (limit > 0) == false {
		limit = 5
	}
	if (page > 0) == false {
		page = 1
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	options := options.Find()
	options.SetSort(bson.D{{"name", 1}})
	options.SetSkip((page - 1) * limit)
	options.SetLimit(limit)

	cursor, err := config.DB.Collection("genres").Find(ctx, bson.M{}, options)
	if err != nil {
		return result, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return result, err
		//panic(err)
	}

	return result, nil
}

func GetTotalGenresCount(query bson.M) (int64, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	fmt.Println(query)
	count, err := config.DB.Collection("genres").CountDocuments(ctx, query, nil)
	return count, err
}

func DeleteGenre(id string) (bool, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return false, err
	}

	books, err := GetBooksByGenreId(id)

	if err != nil {
		return false, err
	}

	if len(books) > 0 {
		return false, errors.New("please delete books before deleting genre")
	}

	_, err = config.DB.Collection("genres").DeleteOne(ctx, bson.M{"_id": oid})

	if err != nil {
		return false, err
	}

	return true, nil
}

// func buildSelQuery(query map[string]string) bson.M {
//     selQuery := make(bson.M, len(query))
//     for k, v := range query {
//         selQuery[k] = v
//     }
//     return selQuery
// }

func UpdateGenre(id string, genre map[string]interface{}) (models.Genre, error) {
	var genreDoc models.Genre
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return genreDoc, err
	}
	query := bson.M{}
	query["_id"] = oid

	updateQuery := bson.M{}
	if genre["name"] != nil {
		updateQuery["name"] = genre["name"].(string)
	}

	if len(updateQuery) == 0 {
		return genreDoc, errors.New("nothing to update")
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	err = config.DB.Collection("genres").FindOneAndUpdate(ctx, query, bson.M{"$set": updateQuery}, &opt).Decode(&genreDoc)

	if err != nil {
		return genreDoc, err
	}

	return genreDoc, nil
}
