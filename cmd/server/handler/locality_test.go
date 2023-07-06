package handler_test

import "testing"

func TestLocalityCreate(t *testing.T) {
	t.Run("Returns 201 if locality is created successfully", func(t *testing.T) {

	})
	t.Run("Returns 409 if locality already exists", func(t *testing.T) {

	})
}

func TestLocalityReport(t *testing.T) {
	t.Run("Returns all localities if id is omitted", func(t *testing.T) {

	})
	t.Run("Returns single locality if id is given", func(t *testing.T) {

	})
	t.Run("Returns 404 if id is not found", func(t *testing.T) {

	})
	t.Run("Returns 201 if no id is given and there is no data", func(t *testing.T) {

	})
}
