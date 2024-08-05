package httpe

import (
	"encoding/json"
	"fmt"
)

type Result struct {
	Code    string `json:"code"`
	Payload any    `json:"payload"`
}

type SuccessResponse struct {
	Result Result `json:"result"`
}

func (er SuccessResponse) MarshalJSON() ([]byte, error) {
	type r SuccessResponse

	var resp r = r(er)

	return json.Marshal(resp)
}

type ErrorResponse struct {
	ErrorStruct struct {
		Code        string `json:"code"`
		Description error  `json:"description"`
	} `json:"error"`
}

func (e ErrorResponse) Error() string {
	return e.ErrorStruct.Description.Error()
}

type Code string

type Message struct {
	Code     Code              `json:"code"`
	Metadata map[string]string `json:"metadata"`
	Payload  any               `json:"payload"`
}

func (m Message) Unwrap() error {
	e, ok := m.Payload.(ErrorResponse)
	if !ok {
		return nil
	}

	var currentErr error = e.ErrorStruct.Description

	for {
		if currentErr == nil {
			return nil
		}

		httpErr, okErr := currentErr.(statusError)
		if okErr {
			return httpErr
		}

		e, ok := currentErr.(interface {
			Unwrap() error
			Error() string
		})
		if !ok {
			return e
		}

		currentErr = e.Unwrap()
	}
}

func (m Message) Error() string {
	e, ok := m.Payload.(ErrorResponse)
	if !ok {
		return ""
	}

	return e.Error()
}

func newMessage(code Code, payload any) (Message, error) {
	return Message{
		Code:     code,
		Metadata: make(map[string]string),
		Payload:  payload,
	}, nil
}

func (m Message) AddMetadata(key, value string) {
	m.Metadata[key] = value
}

func (m Message) MarshalJSON() ([]byte, error) {
	type pack Message

	p := pack(m)

	return json.Marshal(p)
}

func NewErrorMessage(code Code, message string, httpStatusCode int) []byte {
	msg := NewErrorMessageRaw(code, message, httpStatusCode)

	out, _ := msg.MarshalJSON()

	return out
}

type fakeErr struct {
	payload any
	status  int
}

func (e fakeErr) StatusCode() int {
	return e.status
}

func (e fakeErr) Unwrap() error {
	return statusError{
		Err:    e,
		Status: e.status,
	}
}

func (e fakeErr) GetError() error {
	b, _ := e.MarshalJSON()

	return fmt.Errorf("%s", string(b))
}

func (e fakeErr) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.payload)
}

func (e fakeErr) Error() string {
	b, _ := json.Marshal(e.payload)

	return string(b)
}

func NewErrorMessageRaw(code Code, errMessage any, httpStatusCode int) Message {
	var msgErr StatusError

	switch v := errMessage.(type) {
	case string:
		msgErr = statusError{
			Err:    fmt.Errorf("%s", string(v)),
			Status: httpStatusCode,
		}
	case error:
		msgErr = statusError{
			Err:    v,
			Status: httpStatusCode,
		}
	default:
		msgErr = fakeErr{
			payload: errMessage,
			status:  httpStatusCode,
		}
	}

	errMsg := ErrorResponse{
		ErrorStruct: struct {
			Code        string `json:"code"`
			Description error  `json:"description"`
		}{
			Code:        string(code),
			Description: msgErr,
		},
	}

	msg, _ := newMessage(code, errMsg)

	return msg
}

func NewSuccessMessageRaw(code Code, payload any) Message {
	sucessMsg := SuccessResponse{
		Result: struct {
			Code    string `json:"code"`
			Payload any    `json:"payload"`
		}{
			Code:    "success",
			Payload: payload,
		},
	}

	msg, _ := newMessage(code, sucessMsg)

	return msg
}

func NewSuccessMessage(code Code, payload any) []byte {
	sucessMsg := SuccessResponse{
		Result: struct {
			Code    string `json:"code"`
			Payload any    `json:"payload"`
		}{
			Code:    "success",
			Payload: payload,
		},
	}

	msg, _ := newMessage(code, sucessMsg)

	out, _ := msg.MarshalJSON()

	return out
}
