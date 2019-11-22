package httputil

import (
	"errors"
	"net/http"

	"github.com/imtanmoy/authz/utils/sqlutil"

	"github.com/go-chi/render"
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	Code    int                 `json:"code"`    // application-specific error code
	Message string              `json:"message"` // user-level status message
	Errors  map[string][]string `json:"errors"`  // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// HandleError handles sqlerror to http error
func HandleError(err error) render.Renderer {
	var e *sqlutil.SQLError
	sqlerror := sqlutil.GetError(err)
	httpErrCode := 500
	if errors.As(sqlerror, &e) {
		switch e.Code {
		case sqlutil.CodeUniqueViolation:
			httpErrCode = 404
			break
		default:
			httpErrCode = 500
		}
	}

	return &ErrResponse{
		Err:            sqlerror,
		HTTPStatusCode: httpErrCode,
		Message:        sqlerror.Error(),
		Code:           httpErrCode,
		Errors:         make(map[string][]string),
	}
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

func ErrRenderWithCode(err error, code int) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: code,
		Message:        err.Error(),
		Code:           code,
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

func ErrInvalidRequestParamWithMsg(message string) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: 400,
		Message:        message,
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

func ErrUnprocessableEntity() render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: 422,
		Message:        "Request Can not be processed",
		Code:           422,
		Errors:         make(map[string][]string),
	}
}
