package httpe

import (
	"errors"
	"testing"
)

func TestError_Unwrap(t *testing.T) {
	baseError := errors.New("base error")
	baseError2 := errors.New("base error2")

	errHttp500 := statusError{
		Err:    baseError2,
		Status: 500,
	}

	base := Error{
		baseError: baseError,
		httpError: errHttp500,
	}

	if ok := errors.Is(base, errHttp500); !ok {
		t.Error("Want HttpError")
	} else {
		if status := base.httpError.StatusCode(); status != 500 {
			t.Error("Invalid http status")
		}
	}

	if ok := errors.Is(base, baseError); !ok {
		t.Error("Want baseError")
	}

	if ok := errors.Is(base, baseError2); !ok {
		t.Error("Want baseError2")
	}
}
