package controllers

import (
	"context"
	"digitalLibrary/config"
	"digitalLibrary/models"
	"digitalLibrary/utils"
	"errors"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
)

var newCustomLog = utils.NewCustomLog
var convertStrArrToMongoObj = utils.UtilStruct.ConvertStrArrToMongoObj
var isZeroValue = utils.UtilStruct.IsZeroValue
var convertInterfaceArrToStringArr = utils.UtilStruct.ConvertInterfaceArrToStringArr

func GetBook(id string) (models.Book, error) {
	var result models.Book
	//id = "dss"
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return result, err
	}
	//bson.M{"_id": bson.D{{"$gt", 25}}}
	err = config.DB.Collection("books").FindOne(ctx, bson.M{"_id": oid}).Decode(&result)

	if err != nil {
		return result, err
	}

	// if reflect.DeepEqual(result, models.Book{}) == true {
	// 	return models.Book{}, errors.New("Book not found")
	// }

	return result, nil
}

func GetBookPopulated(id string) (map[string]interface{}, error) {
	var result []map[string]interface{}
	var defaultResult map[string]interface{}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return defaultResult, err
	}

	query := []bson.M{
		bson.M{"$match": bson.M{"_id": oid}},
		bson.M{
			"$lookup": bson.M{
				"from":         "authors",
				"localField":   "author",
				"foreignField": "_id",
				"as":           "author",
			},
		},
		bson.M{"$unwind": "$author"},
		bson.M{
			"$lookup": bson.M{
				"from":         "genres",
				"localField":   "genre",
				"foreignField": "_id",
				"as":           "genre",
			},
		},
	}

	cursor, err := config.DB.Collection("books").Aggregate(ctx, query)

	if err != nil {
		return defaultResult, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return defaultResult, err
	}

	if len(result) > 0 {
		return result[0], err
	}

	return defaultResult, nil
}

func GetBooks(limit, page int64) ([]models.Book, error) {
	var result []models.Book
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	if (limit > 0) == false {
		limit = 5
	}
	if (page > 0) == false {
		page = 1
	}
	options := options.Find()
	options.SetSort(bson.D{{"due_back", 1}})
	options.SetSkip((page - 1) * limit)
	options.SetLimit(limit)

	//bson.M{"_id": bson.D{{"$gt", 25}}}
	cursor, err := config.DB.Collection("books").Find(ctx, bson.M{}, options)
	if err != nil {
		return result, err
	}
	if err = cursor.All(ctx, &result); err != nil {
		return result, err
	}

	return result, nil
}

func GetTotalBooksCount(query bson.M) (int64, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	count, err := config.DB.Collection("books").CountDocuments(ctx, query, nil)
	return count, err
}

func GetBooksByAuthorId(id string) ([]models.Book, error) {
	var result []models.Book
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return result, err
	}

	cursor, err := config.DB.Collection("books").Find(ctx, bson.M{"author": oid})

	if err != nil {
		return result, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return result, err
		//panic(err)
	}

	return result, nil
}

func GetBooksByGenreId(id string) ([]models.Book, error) {
	var result []models.Book
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return result, err
	}

	cursor, err := config.DB.Collection("books").Find(ctx, bson.M{"genre": oid})

	if err != nil {
		return result, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return result, err
		//panic(err)
	}

	return result, nil
}

func CreateBook(title, summary, authorIdStr, isbn string, genreArr []string) (interface{}, error) {

	if isZeroValue(title) == true || isZeroValue(summary) == true || isZeroValue(authorIdStr) == true || isZeroValue(genreArr) == true || isZeroValue(isbn) == true {
		return false, errors.New("title, summary, isbn, genre or author missing")
	}

	author, err := primitive.ObjectIDFromHex(authorIdStr)

	if err != nil {
		return nil, err
	}

	genre, err := convertStrArrToMongoObj(genreArr)

	book := models.Book{ID: primitive.NewObjectID(), Title: title, Summary: summary, Author: author, Isbn: isbn, Genre: genre}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	insertResult, err := config.DB.Collection("books").InsertOne(ctx, book)

	if err != nil {
		return nil, err
	}
	newCustomLog("info", fmt.Sprintf("CreateBook book created %s", insertResult), make([]byte, 0), make([]byte, 0), err)
	return insertResult.InsertedID, nil
}

func DeleteBook(id string) (bool, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return false, err
	}

	bookInstances, err := GetBookInstancesByBookId(id)

	if err != nil {
		return false, err
	}

	if len(bookInstances) > 0 {
		return false, errors.New("please delete bookinstances before deleting book")
	}

	_, err = config.DB.Collection("books").DeleteOne(ctx, bson.M{"_id": oid})

	if err != nil {
		return false, err
	}

	return true, nil
}

func UpdateBook(id string, book map[string]interface{}) (models.Book, error) {
	var bookDoc models.Book
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return bookDoc, err
	}
	query := bson.M{}
	query["_id"] = oid
	updateQuery := bson.M{}
	if book["title"] != nil {
		updateQuery["title"] = book["title"].(string)
	}
	if book["summary"] != nil {
		updateQuery["summary"] = book["summary"].(string)
	}
	if book["isbn"] != nil {
		updateQuery["isbn"] = book["isbn"].(string)
	}
	if book["author"] != nil {
		authorIDStr, err := primitive.ObjectIDFromHex(book["author"].(string))
		if err != nil {
			return bookDoc, err
		}
		result, err := GetAuthor(book["author"].(string))
		if reflect.DeepEqual(result, models.Book{}) == true { //check if book exists
			return bookDoc, errors.New("Book not found")
		}
		if err != nil {
			return bookDoc, err
		}
		updateQuery["author"] = authorIDStr
	}
	if book["genre"] != nil {
		genreStringArr, err := convertInterfaceArrToStringArr(book["genre"].([]interface{}))
		if err != nil {
			return bookDoc, err
		}

		genres, err := convertStrArrToMongoObj(genreStringArr)
		if err != nil {
			return bookDoc, err
		}

		//genreCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		count, err := GetTotalGenresCount(bson.M{"_id": bson.M{"$in": genres}})
		//count, err := config.DB.Collection("genres").CountDocuments(genreCtx, bson.M{"_id": bson.M{"$in": genres}}, nil)
		fmt.Println(count, err)
		if err != nil {
			return bookDoc, err
		}
		if count != int64(len(genres)) {
			return bookDoc, errors.New("one or more genre not found")
		}
		updateQuery["genre"] = genres
	}

	if len(updateQuery) == 0 {
		return bookDoc, errors.New("nothing to update")
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	err = config.DB.Collection("books").FindOneAndUpdate(ctx, query, bson.M{"$set": updateQuery}, &opt).Decode(&bookDoc)

	if err != nil {
		return bookDoc, err
	}

	return bookDoc, nil
}

func CreateBookInstance(status, bookIdStr, imprint, dueBackStr string) (interface{}, error) {
	if isZeroValue(status) == true || isZeroValue(bookIdStr) == true || isZeroValue(imprint) == true {
		return nil, errors.New("status, book, isbn, imprint or due_back missing")
	}

	bookIdObject, err := primitive.ObjectIDFromHex(bookIdStr)

	if err != nil {
		return nil, err
	}

	_, err = GetBook(bookIdStr)

	if err != nil {
		return nil, err
	}

	// if reflect.DeepEqual(book, models.Book{}) == true {
	// 	return false, errors.New("book not found, please create book")
	// }

	bookInstance := models.BookInstance{ID: primitive.NewObjectID(), Book: bookIdObject, Status: status, Imprint: imprint}

	var dueBack time.Time
	if isZeroValue(dueBackStr) == false {
		dueBack, err = time.Parse(time.RFC3339, dueBackStr)
		if err != nil {
			return nil, err
		}
		bookInstance.DueBack = dueBack
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	insertResult, err := config.DB.Collection("bookinstances").InsertOne(ctx, bookInstance)
	fmt.Println(bookInstance)
	if err != nil {
		return nil, err
	}
	newCustomLog("info", fmt.Sprintf("CreateBookInstance book instance created %s", insertResult), make([]byte, 0), make([]byte, 0), err)
	return insertResult.InsertedID, nil
}

func GetBookInstance(id string) (models.BookInstance, error) {
	var result models.BookInstance
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return result, err
	}

	err = config.DB.Collection("bookinstances").FindOne(ctx, bson.M{"_id": oid}).Decode(&result)

	if err != nil {
		return result, err
	}

	return result, nil
}

func GetBookInstancePopulated(id string) (map[string]interface{}, error) {
	var result []map[string]interface{}
	var defaultResult map[string]interface{}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return defaultResult, err
	}

	query := []bson.M{
		bson.M{"$match": bson.M{"_id": oid}},
		bson.M{
			"$lookup": bson.M{
				"from":         "books",
				"localField":   "book",
				"foreignField": "_id",
				"as":           "book",
			},
		},
		bson.M{"$unwind": "$book"},
	}

	cursor, err := config.DB.Collection("bookinstances").Aggregate(ctx, query)

	if err != nil {
		return defaultResult, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return defaultResult, err
	}

	if len(result) > 0 {
		return result[0], err
	}

	return defaultResult, nil
}

func GetBookInstancesByBookId(id string) ([]models.BookInstance, error) {
	var result []models.BookInstance
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return result, err
	}

	cursor, err := config.DB.Collection("bookinstances").Find(ctx, bson.M{"book": oid})

	if err != nil {
		return result, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return result, err
		//panic(err)
	}

	return result, nil
}

func GetBookInstances(limit, page int64) ([]models.BookInstance, error) {
	var result []models.BookInstance

	if (limit > 0) == false {
		limit = 5
	}
	if (page > 0) == false {
		page = 1
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	options := options.Find()
	options.SetSort(bson.D{{"due_back", 1}})
	options.SetSkip((page - 1) * limit)
	options.SetLimit(limit)

	cursor, err := config.DB.Collection("bookinstances").Find(ctx, bson.M{}, options)
	if err != nil {
		return result, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return result, err
		//panic(err)
	}

	return result, nil
}

func GetBookInstancesPopulated(limit, page int64) ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	if (limit > 0) == false {
		limit = 5
	}
	if (page > 0) == false {
		page = 1
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	query := []bson.M{
		bson.M{
			"$lookup": bson.M{
				"from":         "books",
				"localField":   "book",
				"foreignField": "_id",
				"as":           "book",
			},
		},
		bson.M{"$unwind": "$book"},
		bson.M{"$sort": bson.M{"due_back": 1}},
		bson.M{"$skip": (page - 1) * limit},
		bson.M{"$limit": limit},
	}

	cursor, err := config.DB.Collection("bookinstances").Aggregate(ctx, query)

	if err != nil {
		return result, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return result, err
	}

	return result, nil
}

func GetTotalBookInstancesCount(query bson.M) (int64, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	count, err := config.DB.Collection("bookinstances").CountDocuments(ctx, query, nil)
	return count, err
}

func DeleteBookInstance(id string) (bool, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return false, err
	}

	_, err = config.DB.Collection("bookinstances").DeleteOne(ctx, bson.M{"_id": oid})

	if err != nil {
		return false, err
	}

	return true, nil
}

func UpdateBookInstance(id string, bookInstance map[string]interface{}) (models.BookInstance, error) {
	var bookInstanceDoc models.BookInstance
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return bookInstanceDoc, err
	}
	query := bson.M{}
	query["_id"] = oid
	updateQuery := bson.M{}
	if bookInstance["status"] != nil {
		updateQuery["status"] = bookInstance["status"].(string)
	}
	if bookInstance["book"] != nil {
		bookIdStr, err := primitive.ObjectIDFromHex(bookInstance["book"].(string))
		if err != nil {
			return bookInstanceDoc, err
		}
		result, err := GetBook(bookInstance["book"].(string))
		if reflect.DeepEqual(result, models.Book{}) == true { //check if book exists
			return bookInstanceDoc, errors.New("Book not found")
		}
		if err != nil {
			return bookInstanceDoc, err
		}
		updateQuery["book"] = bookIdStr
	}
	if bookInstance["imprint"] != nil {
		updateQuery["imprint"] = bookInstance["imprint"].(string)
	}
	if bookInstance["due_back"] != nil {
		dueBackStr, err := time.Parse(time.RFC3339, bookInstance["due_back"].(string))
		if err != nil {
			return bookInstanceDoc, err
		}
		updateQuery["due_back"] = dueBackStr
	}

	if len(updateQuery) == 0 {
		return bookInstanceDoc, errors.New("nothing to update")
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	err = config.DB.Collection("bookinstances").FindOneAndUpdate(ctx, query, bson.M{"$set": updateQuery}, &opt).Decode(&bookInstanceDoc)

	if err != nil {
		return bookInstanceDoc, err
	}

	return bookInstanceDoc, nil
}
