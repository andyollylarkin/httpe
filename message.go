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
	Error struct {
		Code        string `json:"code"`
		Description error  `json:"description"`
	} `json:"error"`
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

	var currentErr error = e.Error.Description

	for {
		if currentErr == nil {
			return nil
		}

		httpErr, okErr := currentErr.(StatusError)
		if okErr {
			return httpErr
		}

		e, ok := currentErr.(interface {
			Unwrap() []error
			Error() string
		})
		if !ok {
			return e
		}

		currentErr = e
	}
}

func (m Message) Error() string {
	type mt Message

	var errM = m

	b, err := json.Marshal(errM)
	if err != nil {
		return ""
	}

	return string(b)
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

func NewErrorMessageRaw(code Code, message string, httpStatusCode int) Message {
	errMsg := ErrorResponse{
		Error: struct {
			Code        string `json:"code"`
			Description error  `json:"description"`
		}{
			Code: string(code),
			Description: StatusError{
				Err:    fmt.Errorf("%s", message),
				Status: httpStatusCode,
			},
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
