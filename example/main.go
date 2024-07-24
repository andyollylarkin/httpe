package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/andyollylarkin/httpe"
	erroradapters "github.com/andyollylarkin/httpe/error_adapters"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{})

	app.Get("/test", func(c *fiber.Ctx) error {
		httpErr := httpe.NewError(errors.New("Not found"), http.StatusNotFound)

		err := erroradapters.FiberResponseWithError(c, httpErr)
		if err != nil {
			log.Println("Write error: ", err.Error())
		}

		return nil
	})

	app.Listen(":9091")
}
