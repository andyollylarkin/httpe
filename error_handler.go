package httpe

import (
	"fmt"
	"net/http"
)

func ResponseWithError(w http.ResponseWriter, err error) error {
	if err == nil {
		return fmt.Errorf("response error cant be nil")
	}

	httpErr, ok := err.(Error)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)

		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			return err
		}
	}

	if httpErr.httpError.err == nil || httpErr.httpError.status == 0 {
		w.WriteHeader(http.StatusInternalServerError)

		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			return err
		}
	}

	w.WriteHeader(httpErr.httpError.Status())
	_, err = w.Write([]byte(httpErr.Error()))
	if err != nil {
		return err
	}

	return nil
}
