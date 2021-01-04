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

func GenreRoutes(route fiber.Router) {
	route.Post("/genre/create", CreateGenreFunc)
	route.Get("/genre/:id", GetGenreFunc)
	route.Get("/genres", GetGenresFunc)
	route.Delete("/genre/:id/delete", DeleteGenreFunc)
	route.Put("/genre/:id/update", UpdateGenreFunc)
}

func CreateGenreFunc(c *fiber.Ctx) error {
	var body map[string]interface{}
	json.Unmarshal(c.Body(), &body)

	if body["name"] == nil {
		sendErrorResponse(c, "genre missing")
		return nil
	}

	genreId, err := controllers.CreateGenre(body["name"].(string))
	if err != nil {
		newCustomLog("error", fmt.Sprintf("CreateGenreFunc call to CreateGenre"), make([]byte, 0), make([]byte, 0), err)
		sendErrorResponse(c, err.Error())
		return nil
	}
	return c.JSON(utils.JsonResponse{Success: true, Msg: "Genre Created", Data: genreId})
}

func GetGenreFunc(c *fiber.Ctx) error {
	author, err := controllers.GetGenre(c.Params("id"))

	if err != nil {
		newCustomLog("error", fmt.Sprintf("GetAuthorFunc call to GetAuthorFunc %s", c.Params("id")), make([]byte, 0), make([]byte, 0), err)
		sendErrorResponse(c, err.Error())
		return nil
	}

	return c.JSON(utils.JsonResponse{Success: true, Data: author})
}

func GetGenresFunc(c *fiber.Ctx) error {
	var limit int64 = 5
	var page int64 = 1
	var err error

	if isZeroValue(c.Query("limit")) != true {
		limit, err = strconv.ParseInt(c.Query("limit"), 10, 64)
		if err != nil {
			newCustomLog("error", fmt.Sprintf("GetGenresFunc parsing limit"), make([]byte, 0), make([]byte, 0), err)
			sendErrorResponse(c, err.Error())
			return nil
		}
	}
	if isZeroValue(c.Query("page")) != true {
		page, err = strconv.ParseInt(c.Query("page"), 10, 64)
		if err != nil {
			newCustomLog("error", fmt.Sprintf("GetGenresFunc parsing page"), make([]byte, 0), make([]byte, 0), err)
			sendErrorResponse(c, err.Error())
			return nil
		}
	}

	authors, err := controllers.GetGenres(limit, page)

	if err != nil {
		newCustomLog("error", fmt.Sprintf("GetGenresFunc call to GetAuthorsFunc"), make([]byte, 0), make([]byte, 0), err)
		sendErrorResponse(c, err.Error())
		return nil
	}

	total, err := controllers.GetTotalBooksCount(bson.M{})
	if err != nil {
		newCustomLog("error", fmt.Sprintf("GetGenresFunc call to GetTotalGenresCount"), make([]byte, 0), make([]byte, 0), err)
		sendErrorResponse(c, err.Error())
		return nil
	}

	return c.JSON(utils.JsonResponse{Success: true, Data: authors, Total: total, Page: page})
}

func DeleteGenreFunc(c *fiber.Ctx) error {
	_, err := controllers.DeleteGenre(c.Params("id"))

	if err != nil {
		newCustomLog("error", fmt.Sprintf("DeleteGenreFunc call to DeleteGenre %s", c.Params("id")), make([]byte, 0), make([]byte, 0), err)
		sendErrorResponse(c, err.Error())
		return nil
	}

	return c.JSON(utils.JsonResponse{Success: true, Msg: "Genre instance Deleted"})
}

func UpdateGenreFunc(c *fiber.Ctx) error {
	var body map[string]interface{}
	json.Unmarshal(c.Body(), &body)

	genre, err := controllers.UpdateGenre(c.Params("id"), body)

	if err != nil {
		newCustomLog("error", fmt.Sprintf("UpdateGenreFunc call to UpdateGenre %s", c.Params("id")), make([]byte, 0), make([]byte, 0), err)
		sendErrorResponse(c, err.Error())
		return nil
	}

	return c.JSON(utils.JsonResponse{Success: true, Msg: "Genre Updated", Data: genre})
}
