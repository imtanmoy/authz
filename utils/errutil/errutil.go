package errutil

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/go-chi/render"
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	Code    int        `json:"code"`    // application-specific error code
	Message string     `json:"message"` // user-level status message
	Errors  url.Values `json:"errors"`  // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if e.Err == nil {
		e.Err = errors.New("Unknown Error")
	}
	if e.Errors == nil {
		e.Errors = make(url.Values)
	}
	if e.Code == 0 {
		e.Code = http.StatusInternalServerError
		e.HTTPStatusCode = http.StatusInternalServerError
	}
	if e.Message == "" {
		e.Message = e.Err.Error()
	}
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		Message:        err.Error(),
		Code:           422,
		Errors:         make(url.Values),
	}
}

func ErrRenderWithCode(err error, code int) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: code,
		Message:        err.Error(),
		Code:           code,
		Errors:         make(url.Values),
	}
}

func ErrInvalidRequest(errors url.Values) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: 400,
		Message:        "Invalid request",
		Code:           400,
		Errors:         errors,
	}
}

func ErrInvalidRequestParam() render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: 400,
		Message:        "Invalid request parameter",
		Code:           400,
		Errors:         make(url.Values),
	}
}

func ErrInvalidRequestParamWithMsg(message string) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: 400,
		Message:        message,
		Code:           400,
		Errors:         make(url.Values),
	}
}

func ErrNotFound(message string) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: 404,
		Message:        message,
		Code:           404,
		Errors:         make(url.Values),
	}
}

func ErrUnprocessableEntity() render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: 422,
		Message:        "Request Can not be processed",
		Code:           422,
		Errors:         make(url.Values),
	}
}

// NewAPIError create new AppError
func NewAPIError(value ...interface{}) *ErrResponse {
	ae := ErrResponse{}
	if len(value) == 0 {
		ae.Err = errors.New("Unknown Error")
	}
	for i, val := range value {
		if i >= 4 {
			break
		}
		switch v := val.(type) {
		case int:
			ae.Code = v
			ae.HTTPStatusCode = v
		case string:
			ae.Message = v
		case error:
			ae.Err = v
		case map[string][]string:
			ae.Errors = v
		}
	}
	if ae.Code == 0 {
		ae.Code = http.StatusInternalServerError
	}
	if ae.Errors == nil {
		ae.Errors = make(url.Values)
	}
	return &ae
}
