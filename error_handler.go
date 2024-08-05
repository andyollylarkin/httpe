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
	httpErrMsg, ok2 := err.(Message)
	if !ok && !ok2 {
		w.WriteHeader(http.StatusInternalServerError)

		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			return err
		}

		return nil
	}

	if ok {
		if httpErr.httpError.GetError() == nil || httpErr.httpError.StatusCode() == 0 {
			w.WriteHeader(http.StatusInternalServerError)

			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				return err
			}
		}

		w.WriteHeader(httpErr.httpError.StatusCode())
		_, err = w.Write([]byte(httpErr.Error()))
		if err != nil {
			return err
		}

		return nil
	} else if ok2 {
		statusErr, ok := httpErrMsg.Unwrap().(StatusError)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)

			_, err := w.Write([]byte(httpErrMsg.Error()))
			if err != nil {
				return err
			}

			return nil
		}

		if statusErr.GetError() == nil || statusErr.StatusCode() == 0 {
			w.WriteHeader(http.StatusInternalServerError)

			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				return err
			}
		}

		jsonErr, err := httpErrMsg.MarshalJSON()
		if err != nil {
			return err
		}

		w.WriteHeader(statusErr.StatusCode())
		_, err = w.Write(jsonErr)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}
