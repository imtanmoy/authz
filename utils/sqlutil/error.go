package sqlutil

import (
	"errors"
	"fmt"

	"github.com/go-pg/pg/v9"
)

var (
	// ErrRecordNotFound custom error for sql not found
	ErrRecordNotFound = errors.New("record not found")
	// ErrUnknown custom error for sql unknown err
	ErrUnknownError = errors.New("somthing went wrong")

	ErrAlreadyExists = errors.New("record already exists")

	ErrMultiRows = errors.New("more than one record found")
)

const (
	CodeNumericValueOutOfRange    = "22003"
	CodeInvalidTextRepresentation = "22P02"
	CodeNotNullViolation          = "23502"
	CodeForeignKeyViolation       = "23503"
	CodeUniqueViolation           = "23505"
	CodeCheckViolation            = "23514"
	CodeLockNotAvailable          = "55P03"
)

// GetError utility function for converting sql error to custom error
func GetError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pg.ErrNoRows) {
		return ErrRecordNotFound
	}
	if errors.Is(err, pg.ErrMultiRows) {
		return ErrMultiRows
	}
	if pgErr, ok := err.(pg.Error); ok {
		fmt.Println(pgErr.Field('C'))
		fmt.Println(pgErr.Field('M'))
		fmt.Println(pgErr.Field('D'))
		fmt.Println(pgErr.Field('t'))
		fmt.Println(pgErr.Field('c'))
		fmt.Println(pgErr.Field('H'))
		switch pgErr.Field('C') {
		case CodeUniqueViolation:
			return ErrAlreadyExists
		default:
			return ErrUnknownError
		}

	}
	return ErrUnknownError
}
