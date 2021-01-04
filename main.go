package main

import (
	"digitalLibrary/config"
	"digitalLibrary/routes"
	"digitalLibrary/utils"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
	//establish mongo connection
	config.ConnectToDB()
	// utils.TestLog()
}

// type CustomLogInterfaceMock interface {
// 	WriteDevLogs()
// 	WriteLog()
// }

// type CustomLog struct {
// }

// type CustomLogFactory struct{}

// func (c *CustomLog) WriteLog() {
// 	if os.Getenv("GO_ENV") == "production" { //production
// 		c.WriteDevLogs()
// 	} else { //development
// 		c.WriteDevLogs()
// 	}
// }

// func (c *CustomLog) WriteDevLogs() {
// 	fmt.Println("dddddd")
// }

// func (c *CustomLogFactory) NewCustomLog(Level string, Msg string, LogPayload []byte, ExtraInfo []byte, Err error) utils.CustomLogInterface {
// 	return &CustomLog{}
// }

func main() {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())

	routes.SetUpRoutes(app)

	// utils.CustomLogFactoryStruct = &CustomLogFactory{}
	utils.TestLog()
	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.SendString("Hello, World 👋!")
	// })

	host := ":" + os.Getenv("HTTP_PORT")
	log.Fatal(app.Listen(host))
}
