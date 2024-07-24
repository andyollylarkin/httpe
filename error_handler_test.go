package httpe

import (
	"errors"
	"net/http"
	"testing"
)

type MockWriter struct {
	resp   []byte
	status int
}

func (mw MockWriter) Header() http.Header {
	return http.Header{}
}

func (mw *MockWriter) Write(body []byte) (int, error) {
	mw.resp = body

	return len(body), nil
}

func (mw *MockWriter) WriteHeader(statusCode int) {
	mw.status = statusCode
}

func TestResponseWithError(t *testing.T) {
	err := NewError(errors.New("object not found"), http.StatusNotFound)

	rw := &MockWriter{
		resp: []byte{},
	}

	errWrite := ResponseWithError(rw, err)
	if errWrite != nil {
		t.Fatal(errWrite)
	}

	if rw.status != 404 {
		t.Error("Status not written")
	}

	if string(rw.resp) != "object not found" {
		t.Error("Body not written")
	}
}
