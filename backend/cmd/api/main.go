package main

import (
	"fmt"
	"os"
	"projeto_mecanica/internal/config"

	"github.com/gofiber/fiber/v3"
)

func main() {
	cfg, err := config.LoadENV()
	if err != nil {
		fmt.Printf("Error to load config: %v", err)
		os.Exit(1)
	}

	db, err := config.InitDB(cfg)
	if err != nil {
		fmt.Printf("Error to load database: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Conected to database: %v", db)

	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Welcome do Mechanics.")
	})

	app.Listen(":8080")
}
