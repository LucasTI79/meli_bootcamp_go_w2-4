package product

import "fmt"

type ErrInvalidProductCode struct {
	Code string
}

type ErrNotFound struct {
	ID int
}

func (e ErrInvalidProductCode) Error() string {
	return fmt.Sprintf("invalid product code: %s is not unique", e.Code)
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("product with ID %d not found", e.ID)
}
