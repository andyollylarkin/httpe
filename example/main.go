package main

import (
	"net/http"

	"github.com/andyollylarkin/httpe"
	erroradapters "github.com/andyollylarkin/httpe/error_adapters"
	"github.com/gofiber/fiber/v2"
)

func main() {
	m := httpe.NewErrorMessageRaw(httpe.Code("BAD"), struct {
		Desc string `json:"desc"`
		Out  string `json:"out"`
	}{
		Desc: "Hello",
		Out:  "World",
	}, http.StatusBadRequest)

	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		return erroradapters.FiberResponseWithError(c, m)
	})

	app.Listen(":10001")
}
