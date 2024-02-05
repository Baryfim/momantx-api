package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sixfwa/fiber-api/database"
	"github.com/sixfwa/fiber-api/middlewares"
	"github.com/sixfwa/fiber-api/routes"
)

func welcome(c *fiber.Ctx) error {
	return c.SendString("Welcome to my awesome API")
}

func setupRoutes(app *fiber.App) {
	// welcome
	app.Get("/api", welcome)

	// Year endpoints
	app.Get("/api/years/:id", routes.GetYear)
	app.Get("/api/years", routes.GetYears)

	app.Post("/api/years", middlewares.CheckAdminIsValid, routes.CreateYear)
	app.Put("/api/years/:id", middlewares.CheckAdminIsValid, routes.UpdateYear)
	app.Delete("/api/years/:id", middlewares.CheckAdminIsValid, routes.DeleteYear)

	// Item endpoints
	app.Get("/api/items", routes.GetItems)
	app.Get("/api/items/:id", routes.GetItem)

	app.Post("/api/items", middlewares.CheckAdminIsValid, routes.CreateItem)
	app.Put("/api/items/:id", middlewares.CheckAdminIsValid, routes.UpdateItem)
	app.Delete("/api/items/:id", middlewares.CheckAdminIsValid, routes.DeleteItem)

	// Test endpoints
	app.Get("/api/tests/:id", routes.GetTest)
	app.Get("/api/tests", routes.GetTests)

	app.Post("/api/tests", middlewares.CheckAdminIsValid, routes.CreateTest)
	app.Put("/api/tests/:id", middlewares.CheckAdminIsValid, routes.UpdateTest)
	app.Delete("/api/tests/:id", middlewares.CheckAdminIsValid, routes.DeleteTest)

	// Questions endpoints
	app.Get("/api/questions/:id", routes.GetQuestion)
	app.Get("/api/questions", routes.GetQuestions)
	app.Post("/api/questions", middlewares.CheckAdminIsValid, routes.CreateQuestion)
	app.Put("/api/questions/:id", middlewares.CheckAdminIsValid, routes.UpdateQuestion)
	app.Delete("/api/questions/:id", middlewares.CheckAdminIsValid, routes.DeleteQuestion)
}

func main() {
	database.ConnectDb()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept, User-Agent, Authorization",
		AllowCredentials: true,
	}))

	setupRoutes(app)

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
