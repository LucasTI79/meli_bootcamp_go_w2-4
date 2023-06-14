package product

import "fmt"

type ErrInvalidProductCode struct {
	Code string
}

type ErrNotFound struct {
	ID int
}

type ErrGeneric struct {
	message string
}

func NewErrGeneric(message string) *ErrGeneric {
	return &ErrGeneric{message}
}

func (e ErrGeneric) Error() string {
	return e.message
}

func NewErrInvalidProductCode(code string) *ErrInvalidProductCode {
	return &ErrInvalidProductCode{code}
}

func (e ErrInvalidProductCode) Error() string {
	return fmt.Sprintf("invalid product code: %s is not unique", e.Code)
}

func NewErrNotFound(id int) *ErrNotFound {
	return &ErrNotFound{id}
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("product with ID %d not found", e.ID)
}
