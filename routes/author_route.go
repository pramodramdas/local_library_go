package routes

import (
	"digitalLibrary/controllers"
	"digitalLibrary/utils"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// var newCustomLog = utils.NewCustomLog
// var sendErrorResponse = utils.UtilStruct.SendErrorResponse
// var isZeroValue = utils.UtilStruct.IsZeroValue
// var convertInterfaceArrToStringArr = utils.UtilStruct.ConvertInterfaceArrToStringArr

func AuthorRoutes(route fiber.Router) {
	route.Post("/author/create", CreateAuthorFunc)
	route.Get("/author/:id", GetAuthorFunc)
	route.Get("/authors", GetAuthorsFunc)
	route.Delete("/author/:id/delete", DeleteAuthorFunc)
	route.Put("/author/:id/update", UpdateAuthorFunc)
}

func CreateAuthorFunc(c *fiber.Ctx) error {
	var body map[string]interface{}
	json.Unmarshal(c.Body(), &body)

	if body["first_name"] == nil || body["family_name"] == nil || body["date_of_birth"] == nil {
		sendErrorResponse(c, "first_name, family_name or date_of_birth missing")
		return nil
	}

	if body["date_of_death"] == nil {
		body["date_of_death"] = ""
	}
	authorId, err := controllers.CreateAuthor(body["first_name"].(string), body["family_name"].(string), body["date_of_birth"].(string), body["date_of_death"].(string))
	if err != nil {
		newCustomLog("error", fmt.Sprintf("CreateAuthorFunc call to CreateAuthor"), make([]byte, 0), make([]byte, 0), err)
		sendErrorResponse(c, err.Error())
		return nil
	}
	return c.JSON(utils.JsonResponse{Success: true, Msg: "Author Created", Data: authorId})
}

func GetAuthorFunc(c *fiber.Ctx) error {
	author, err := controllers.GetAuthor(c.Params("id"))

	if err != nil {
		newCustomLog("error", fmt.Sprintf("GetAuthorFunc call to GetAuthor %s", c.Params("id")), make([]byte, 0), make([]byte, 0), err)
		sendErrorResponse(c, err.Error())
		return nil
	}

	return c.JSON(utils.JsonResponse{Success: true, Data: author})
}

func GetAuthorsFunc(c *fiber.Ctx) error {
	var limit int64 = 5
	var page int64 = 1
	var err error

	if isZeroValue(c.Query("limit")) != true {
		limit, err = strconv.ParseInt(c.Query("limit"), 10, 64)
		if err != nil {
			newCustomLog("error", fmt.Sprintf("GetAuthorsFunc parsing limit"), make([]byte, 0), make([]byte, 0), err)
			sendErrorResponse(c, err.Error())
			return nil
		}
	}
	if isZeroValue(c.Query("page")) != true {
		page, err = strconv.ParseInt(c.Query("page"), 10, 64)
		if err != nil {
			newCustomLog("error", fmt.Sprintf("GetAuthorsFunc parsing page"), make([]byte, 0), make([]byte, 0), err)
			sendErrorResponse(c, err.Error())
			return nil
		}
	}

	authors, err := controllers.GetAuthors(limit, page)

	if err != nil {
		newCustomLog("error", fmt.Sprintf("GetAuthorsFunc call to GetAuthors"), make([]byte, 0), make([]byte, 0), err)
		sendErrorResponse(c, err.Error())
		return nil
	}

	total, err := controllers.GetTotalAuthorsCount(bson.M{})
	if err != nil {
		newCustomLog("error", fmt.Sprintf("GetAuthorsFunc call to GetTotalAuthorsCount"), make([]byte, 0), make([]byte, 0), err)
		sendErrorResponse(c, err.Error())
		return nil
	}

	return c.JSON(utils.JsonResponse{Success: true, Data: authors, Total: total, Page: page})
}

func DeleteAuthorFunc(c *fiber.Ctx) error {
	_, err := controllers.DeleteAuthor(c.Params("id"))

	if err != nil {
		newCustomLog("error", fmt.Sprintf("DeleteAuthorFunc call to DeleteAuthor %s", c.Params("id")), make([]byte, 0), make([]byte, 0), err)
		sendErrorResponse(c, err.Error())
		return nil
	}

	return c.JSON(utils.JsonResponse{Success: true, Msg: "Author Deleted"})
}

func UpdateAuthorFunc(c *fiber.Ctx) error {
	var body map[string]interface{}
	json.Unmarshal(c.Body(), &body)

	author, err := controllers.UpdateAuthor(c.Params("id"), body)

	if err != nil {
		newCustomLog("error", fmt.Sprintf("UpdateAuthorFunc call to UpdateAuthor %s", c.Params("id")), make([]byte, 0), make([]byte, 0), err)
		sendErrorResponse(c, err.Error())
		return nil
	}

	return c.JSON(utils.JsonResponse{Success: true, Msg: "Author Updated", Data: author})
}
