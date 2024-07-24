package main

import (
	"log"
	"net/http"

	"github.com/andyollylarkin/httpe"
	erroradapters "github.com/andyollylarkin/httpe/error_adapters"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{})

	app.Get("/test", func(c *fiber.Ctx) error {
		msgErr := httpe.NewErrorMessageRaw(httpe.Code("TEST"), "some error", http.StatusBadRequest)

		err := erroradapters.FiberResponseWithError(c, msgErr)
		if err != nil {
			log.Println("Write error: ", err.Error())
		}

		return nil
	})

	app.Listen(":9091")
}
