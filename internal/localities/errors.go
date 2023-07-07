package localities

import (
	"fmt"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
)

type ErrInvalidLocality struct {
	Locality domain.Locality
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

func NewErrInvalidLocality(loc domain.Locality) *ErrInvalidLocality {
	return &ErrInvalidLocality{loc}
}

func (e ErrInvalidLocality) Error() string {
	return fmt.Sprintf("locality already exists: %s, %s, %s", e.Locality.Name, e.Locality.Province, e.Locality.Country)
}

func NewErrNotFound(id int) *ErrNotFound {
	return &ErrNotFound{id}
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("locality with ID %d not found", e.ID)
}
