package routes

import "github.com/gofiber/fiber/v2"

func SetUpRoutes(app *fiber.App) {
	catalog := app.Group("/catalog")

	BookRoutes(catalog)
	AuthorRoutes(catalog)
	GenreRoutes(catalog)
	IndexRoutes(catalog)
}
