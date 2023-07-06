package localities_test

import (
	"context"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/localities"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	t.Run("Creates valid locality", func(t *testing.T) {
		t.Skip()
	})
	t.Run("Doesn't create duplicate locality", func(t *testing.T) {
		t.Skip()
	})
}

func TestReport(t *testing.T) {
	t.Run("Returns all localities when id is omitted", func(t *testing.T) {
		t.Skip()
	})
	t.Run("Returns specific locality if id is provided", func(t *testing.T) {
		t.Skip()
	})
	t.Run("Returns NotFound if id doesn't exist", func(t *testing.T) {
		t.Skip()
	})
	t.Run("Doesn't return NotFound if id is omitted but no data exists", func(t *testing.T) {
		t.Skip()
	})
}

