package handler_test

import (
	"testing"
)

func TestProductCreate(t *testing.T) {
	t.Run("Returns 201 when creation succeds", func(t *testing.T) {
		t.Skip()
	})
	t.Run("Does not fail when fields are zero", func(t *testing.T) {
		t.Skip()
	})
	t.Run("Returns 422 when required fields are omitted", func(t *testing.T) {
		t.Skip()
	})
	t.Run("Returns 409 when product code is not unique", func(t *testing.T) {
		t.Skip()
	})
}

func TestProductRead(t *testing.T) {
	t.Run("Returns all products on GetAll", func(t *testing.T) {
		t.Skip()
	})
	t.Run("Returns 204 when GetAll returns no products", func(t *testing.T) {
		t.Skip()
	})
	t.Run("Returns existing product on Get by ID", func(t *testing.T) {
		t.Skip()
	})
	t.Run("Returns 400 when ID is not an int", func(t *testing.T) {
		t.Skip()
	})
}

func TestProductUpdate(t *testing.T) {
	t.Run("Returns 200 when update succeds", func(t *testing.T) {
		t.Skip()
	})
	t.Run("Does not fail when updated value is zero", func(t *testing.T) {
		t.Skip()
	})
	t.Run("Returns 404 when ID is not found", func(t *testing.T) {
		t.Skip()
	})
	t.Run("Returns 409 when updated code exists", func(t *testing.T) {
		t.Skip()
	})
}

func TestProductDelete(t *testing.T) {

	t.Run("Returns 200 when delete succeds", func(t *testing.T) {
		t.Skip()
	})
	t.Run("Returns 404 when ID is not found", func(t *testing.T) {
		t.Skip()
	})
}
