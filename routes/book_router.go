package routes

import (
	"digitalLibrary/controllers"
	"digitalLibrary/utils"
	"encoding/json"
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gofiber/fiber/v2"
)

var newCustomLog = utils.NewCustomLog
var sendErrorResponse = utils.UtilStruct.SendErrorResponse
var isZeroValue = utils.UtilStruct.IsZeroValue
var convertInterfaceArrToStringArr = utils.UtilStruct.ConvertInterfaceArrToStringArr

func BookRoutes(route fiber.Router) {
	route.Get("/book/:id", GetBookFunc)
	route.Get("/books", GetBooksFunc)
	route.Post("/book/create", CreateBookFunc)
	route.Delete("/book/:id/delete", DeleteBookFunc)
	route.Put("/book/:id/update", UpdateBookFunc)

	route.Post("/bookinstance/create", CreateBookInstanceFunc)
	route.Get("/bookinstance/:id", GetBookInstanceFunc)
	route.Get("/bookinstances", GetBookInstancesFunc)
	route.Get("/bookinstancesbybookid/:id", GetBookInstancesByBookIdFunc)
	route.Get("/booksbygenreid/:id", GetBooksByGenreIdFunc)
	route.Get("/booksbyauthorid/:id", GetBooksByAuthorIdFunc)
	route.Delete("/bookinstance/:id/delete", DeleteBookInstanceFunc)
	route.Put("/bookinstance/:id/update", UpdateBookInstanceFunc)
}

func GetBookFunc(c *fiber.Ctx) error {
	var book interface{}
	var err error

	if c.Query("populate") == "yes" {
		book, err = controllers.GetBookPopulated(c.Params("id"))
	} else {
		book, err = controllers.GetBook(c.Params("id"))
	}

	if err != nil {
		newCustomLog("error", fmt.Sprintf("getBookFunc call to GetBook %s", c.Params("id")), make([]byte, 0), make([]byte, 0), err)
		sendErrorResponse(c, err.Error())
		return nil
	}

	return c.JSON(utils.JsonResponse{Success: true, Data: book})
}

func GetBooksFunc(c *fiber.Ctx) error {
	var limit int64 = 5
	var page int64 = 1
	var err error

	if isZeroValue(c.Query("limit")) != true {
		limit, err = strconv.ParseInt(c.Query("limit"), 10, 64)
		if err != nil {
			newCustomLog("error", fmt.Sprintf("GetBooksFunc parsing limit"), make([]byte, 0), make([]byte, 0), err)
			sendErrorResponse(c, err.Error())
			return nil
		}
	}
	if isZeroValue(c.Query("page")) != true {
		page, err = strconv.ParseInt(c.Query("page"), 10, 64)
		if err != nil {
			newCustomLog("error", fmt.Sprintf("GetBooksFunc parsing page"), make([]byte, 0), make([]byte, 0), err)
			sendErrorResponse(c, err.Error())
			return nil
		}
	}

	books, err := controllers.GetBooks(limit, page)

	if err != nil {
		newCustomLog("error", fmt.Sprintf("getBooksFunc call to GetBooks"), make([]byte, 0), make([]byte, 0), err)
		sendErrorResponse(c, err.Error())
		return nil
	}

	total, err := controllers.GetTotalBooksCount(bson.M{})
	if err != nil {
		newCustomLog("error", fmt.Sprintf("getBooksFunc call to GetTotalBooksCount"), make([]byte, 0), make([]byte, 0), err)
		sendErrorResponse(c, err.Error())
		return nil
	}

	return c.JSON(utils.JsonResponse{Success: true, Data: books, Total: total, Page: page})
	//return c.Send(data)
}

func GetBooksByGenreIdFunc(c *fiber.Ctx) error {
	books, err := controllers.GetBooksByGenreId(c.Params("id"))

	if err != nil {
		newCustomLog("error", fmt.Sprintf("GetBooksByGenreIdFunc call to GetBooksByGenreId %s", c.Params("id")), make([]byte, 0), make([]byte, 0), err)
		sendErrorResponse(c, err.Error())
		return nil
	}

	return c.JSON(utils.JsonResponse{Success: true, Data: books})
}

func GetBooksByAuthorIdFunc(c *fiber.Ctx) error {
	books, err := controllers.GetBooksByAuthorId(c.Params("id"))

	if err != nil {
		newCustomLog("error", fmt.Sprintf("GetBooksByAuthorIdFunc call to GetBooksByAuthorId %s", c.Params("id")), make([]byte, 0), make([]byte, 0), err)
		sendErrorResponse(c, err.Error())
		return nil
	}

	return c.JSON(utils.JsonResponse{Success: true, Data: books})
}

func CreateBookFunc(c *fiber.Ctx) error {
	//var book models.Book
	var body map[string]interface{}
	json.Unmarshal(c.Body(), &body)

	if body["title"] == nil || body["summary"] == nil || body["author"] == nil || body["genre"] == nil || body["isbn"] == nil {
		sendErrorResponse(c, "title, summary, isbn, genre or author missing")
		return nil
	}

	genreStringArr, err := convertInterfaceArrToStringArr(body["genre"].([]interface{}))

	bookId, err := controllers.CreateBook(body["title"].(string), body["summary"].(string), body["author"].(string), body["isbn"].(string), genreStringArr)
	if err != nil {
		newCustomLog("error", fmt.Sprintf("CreateBookFunc call to CreateBook"), make([]byte, 0), make([]byte, 0), err)
		sendErrorResponse(c, err.Error())
		return nil
	}
	return c.JSON(utils.JsonResponse{Success: true, Msg: "Book Created", Data: bookId})
}

func DeleteBookFunc(c *fiber.Ctx) error {
	_, err := controllers.DeleteBook(c.Params("id"))

	if err != nil {
		newCustomLog("error", fmt.Sprintf("DeleteBookFunc call to DeleteBook %s", c.Params("id")), make([]byte, 0), make([]byte, 0), err)
		sendErrorResponse(c, err.Error())
		return nil
	}

	return c.JSON(utils.JsonResponse{Success: true, Msg: "Book Deleted"})
}

func UpdateBookFunc(c *fiber.Ctx) error {
	var body map[string]interface{}
	json.Unmarshal(c.Body(), &body)

	book, err := controllers.UpdateBook(c.Params("id"), body)

	if err != nil {
		newCustomLog("error", fmt.Sprintf("UpdateBookFunc call to UpdateBook %s", c.Params("id")), make([]byte, 0), make([]byte, 0), err)
		sendErrorResponse(c, err.Error())
		return nil
	}

	return c.JSON(utils.JsonResponse{Success: true, Msg: "Book Updated", Data: book})
}

func CreateBookInstanceFunc(c *fiber.Ctx) error {
	var body map[string]interface{}
	json.Unmarshal(c.Body(), &body)

	if body["status"] == nil || body["book"] == nil || body["imprint"] == nil || body["due_back"] == nil {
		sendErrorResponse(c, "status, book, imprint, isbn, due_back missing")
		return nil
	}

	bookInstanceId, err := controllers.CreateBookInstance(body["status"].(string), body["book"].(string), body["imprint"].(string), body["due_back"].(string))
	if err != nil {
		newCustomLog("error", fmt.Sprintf("CreateBookInstanceFunc call to CreateBookInstance"), make([]byte, 0), make([]byte, 0), err)
		sendErrorResponse(c, err.Error())
		return nil
	}
	return c.JSON(utils.JsonResponse{Success: true, Msg: "Book instance Created", Data: bookInstanceId})
}

func GetBookInstanceFunc(c *fiber.Ctx) error {
	var bookInstance interface{}
	var err error
	if c.Query("populate") == "yes" {
		bookInstance, err = controllers.GetBookInstancePopulated(c.Params("id"))
	} else {
		bookInstance, err = controllers.GetBookInstance(c.Params("id"))
	}

	if err != nil {
		newCustomLog("error", fmt.Sprintf("GetBookInstanceFunc call to GetBookInstance %s", c.Params("id")), make([]byte, 0), make([]byte, 0), err)
		sendErrorResponse(c, err.Error())
		return nil
	}

	return c.JSON(utils.JsonResponse{Success: true, Data: bookInstance})
}

func GetBookInstancesByBookIdFunc(c *fiber.Ctx) error {
	bookInstances, err := controllers.GetBookInstancesByBookId(c.Params("id"))
	if err != nil {
		newCustomLog("error", fmt.Sprintf("GetBookInstancesByBookIdFunc call to GetBookInstancesByBookId %s", c.Params("id")), make([]byte, 0), make([]byte, 0), err)
		sendErrorResponse(c, err.Error())
		return nil
	}

	return c.JSON(utils.JsonResponse{Success: true, Data: bookInstances})
}

func GetBookInstancesFunc(c *fiber.Ctx) error {
	var limit int64 = 5
	var page int64 = 1
	var err error

	if isZeroValue(c.Query("limit")) != true {
		limit, err = strconv.ParseInt(c.Query("limit"), 10, 64)
		if err != nil {
			newCustomLog("error", fmt.Sprintf("GetBookInstancesFunc parsing limit"), make([]byte, 0), make([]byte, 0), err)
			sendErrorResponse(c, err.Error())
			return nil
		}
	}
	if isZeroValue(c.Query("page")) != true {
		page, err = strconv.ParseInt(c.Query("page"), 10, 64)
		if err != nil {
			newCustomLog("error", fmt.Sprintf("GetBookInstancesFunc parsing page"), make([]byte, 0), make([]byte, 0), err)
			sendErrorResponse(c, err.Error())
			return nil
		}
	}

	var bookInstances interface{}
	if c.Query("populate") == "yes" {
		bookInstances, err = controllers.GetBookInstancesPopulated(limit, page)
	} else {
		bookInstances, err = controllers.GetBookInstances(limit, page)
	}

	if err != nil {
		newCustomLog("error", fmt.Sprintf("GetBookInstancesFunc call to GetBookInstances"), make([]byte, 0), make([]byte, 0), err)
		sendErrorResponse(c, err.Error())
		return nil
	}

	total, err := controllers.GetTotalBookInstancesCount(bson.M{})
	if err != nil {
		newCustomLog("error", fmt.Sprintf("GetBookInstancesFunc call to GetTotalBookInstancesCount"), make([]byte, 0), make([]byte, 0), err)
		sendErrorResponse(c, err.Error())
		return nil
	}

	return c.JSON(utils.JsonResponse{Success: true, Data: bookInstances, Total: total, Page: page})
}

func DeleteBookInstanceFunc(c *fiber.Ctx) error {
	_, err := controllers.DeleteBookInstance(c.Params("id"))

	if err != nil {
		newCustomLog("error", fmt.Sprintf("DeleteBookInstanceFunc call to DeleteBookInstance %s", c.Params("id")), make([]byte, 0), make([]byte, 0), err)
		sendErrorResponse(c, err.Error())
		return nil
	}

	return c.JSON(utils.JsonResponse{Success: true, Msg: "Book instance Deleted"})
}

func UpdateBookInstanceFunc(c *fiber.Ctx) error {
	var body map[string]interface{}
	json.Unmarshal(c.Body(), &body)

	bookInstance, err := controllers.UpdateBookInstance(c.Params("id"), body)

	if err != nil {
		newCustomLog("error", fmt.Sprintf("UpdateBookInstanceFunc call to UpdateBookInstance %s", c.Params("id")), make([]byte, 0), make([]byte, 0), err)
		sendErrorResponse(c, err.Error())
		return nil
	}

	return c.JSON(utils.JsonResponse{Success: true, Msg: "Book Instance Updated", Data: bookInstance})
}
