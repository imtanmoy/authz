package httputil

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/imtanmoy/authz/utils/sqlutil"

	"github.com/go-chi/render"
)

// ErrResponse construct http error response
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	Code    int        `json:"code"`    // application-specific error code
	Message string     `json:"message"` // user-level status message
	Errors  url.Values `json:"errors"`  // application-level error message, for debugging
}

// Render ErrResponse render method for chi
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// HandleError handles sqlerror to http error
// func HandleError(err error) *ErrResponse {
// 	var e *sqlutil.SQLError
// 	sqlerror := sqlutil.GetError(err)
// 	httpErrCode := 500
// 	if errors.As(sqlerror, &e) {
// 		switch e.Code {
// 		case sqlutil.CodeUniqueViolation:
// 			httpErrCode = 404
// 			break
// 		default:
// 			httpErrCode = 500
// 		}
// 	}

// 	return &ErrResponse{
// 		Err:            sqlerror,
// 		HTTPStatusCode: httpErrCode,
// 		Message:        sqlerror.Error(),
// 		Code:           httpErrCode,
// 		Errors:         make(map[string][]string),
// 	}
// }

func handleSQLError(err error) *ErrResponse {
	var e *sqlutil.SQLError
	httpErrCode := 500
	if errors.As(err, &e) {
		switch e.Code {
		case sqlutil.CodeUniqueViolation:
			httpErrCode = 404
			break
		default:
			httpErrCode = 500
		}
	}

	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: httpErrCode,
		Message:        err.Error(),
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
			break
		case string:
			ae.Message = v
			break
		case error:
			sqlerror := sqlutil.GetError(v)
			if sqlerror != nil && len(value) == 1 {
				ae = *handleSQLError(sqlerror)
				break
			}
			ae.Err = v
		case map[string][]string:
			ae.Errors = v
			break
		}
	}
	if ae.Code == 0 {
		ae.Code = http.StatusInternalServerError
		ae.HTTPStatusCode = http.StatusInternalServerError
	}
	if ae.Errors == nil {
		ae.Errors = make(url.Values)
	}
	if ae.Err == nil {
		ae.Err = errors.New("Unknown Error")
	}
	if ae.Message == "" {
		ae.Message = ae.Err.Error()
	}
	return &ae
}
