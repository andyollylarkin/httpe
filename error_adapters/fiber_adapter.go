package erroradapters

import (
	"encoding/json"
	"net/http"

	"github.com/andyollylarkin/httpe"
	"github.com/gofiber/fiber/v2"
)

type fiberWriter struct {
	out    []byte
	status int
	header http.Header
}

func (fw *fiberWriter) Header() http.Header {
	return fw.header
}

func (fw *fiberWriter) Write(b []byte) (int, error) {
	fw.out = b

	return len(b), nil
}

func (fw *fiberWriter) WriteHeader(statusCode int) {
	fw.status = statusCode
}

func FiberResponseWithError(c *fiber.Ctx, err error) error {
	resp := &fiberWriter{
		out: make([]byte, 0),
	}

	_, ok := err.(json.Marshaler)

	if ok {
		c.Set("Content-Type", fiber.MIMEApplicationJSON)
	}

	writeErr := httpe.ResponseWithError(resp, err)
	if writeErr != nil {
		return writeErr
	}

	c.Status(resp.status)

	return c.Send(resp.out)
}
