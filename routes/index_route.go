package routes

import (
	"digitalLibrary/controllers"
	"digitalLibrary/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func IndexRoutes(route fiber.Router) {
	route.Get("*", GetAllCountFunc)
}

func GetAllCountFunc(c *fiber.Ctx) error {
	allCounts, errs := controllers.GetAllCount()
	if len(errs) > 0 {
		fmt.Println(errs)
	}
	return c.JSON(utils.JsonResponse{Success: true, Data: allCounts})
}
