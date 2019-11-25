package httputil

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/imtanmoy/authz/utils/sqlutil"

	"github.com/go-chi/render"
)

var ErrInternalServerError = errors.New("internal server error")

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
		Errors:         make(url.Values),
	}
}

// NewAPIError create new AppError
func NewAPIError(value ...interface{}) *ErrResponse {
	ae := ErrResponse{}

	if len(value) == 0 {
		ae.Err = ErrInternalServerError
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
		case url.Values:
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
		ae.Err = ErrInternalServerError
	}
	if ae.Message == "" {
		if errors.Is(ae.Err, ErrInternalServerError) {
			ae.Message = "Oops! Something went wrong"
		} else {
			ae.Message = ae.Err.Error()
		}
	}
	return &ae
}
