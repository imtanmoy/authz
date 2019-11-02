package errutil

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	Message string              `json:"message"`          // user-level status message
	Code    int                 `json:"code"`             // application-specific error code
	Errors  map[string][]string `json:"errors"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		Message:        err.Error(),
		Code:           422,
		Errors:         make(map[string][]string),
	}
}

func ErrInvalidRequest(errors map[string][]string) render.Renderer {
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
		Errors:         make(map[string][]string),
	}
}

func ErrNotFound(message string) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: 404,
		Message:        message,
		Code:           404,
		Errors:         make(map[string][]string),
	}
}
