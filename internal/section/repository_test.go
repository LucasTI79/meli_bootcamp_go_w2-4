package section_test

import (
	"context"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/section"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestRepoCreate(t *testing.T) {
	t.Run("Creates valid section", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := section.NewRepository(db)
		sec := getTestSection()

		id, err := repo.Save(context.TODO(), sec)

		assert.NoError(t, err)

		var receivedNumber int
		row := db.QueryRow(`SELECT section_number FROM sections WHERE id = ?;`, id)
		row.Scan(&receivedNumber)

		assert.Equal(t, sec.SectionNumber, receivedNumber)
	})
}

func TestRepoExists(t *testing.T) {
	t.Run("Finds existing section", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := section.NewRepository(db)
		expected := getTestSection()

		repo.Save(context.TODO(), expected)

		exists := repo.Exists(context.TODO(), expected.SectionNumber)

		assert.True(t, exists)
	})
}

func TestRepoGet(t *testing.T) {
	t.Run("Gets section", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := section.NewRepository(db)
		expected := getTestSection()

		id, _ := repo.Save(context.TODO(), expected)

		received, err := repo.Get(context.TODO(), id)
		assert.NoError(t, err)

		assert.Equal(t, expected.SectionNumber, received.SectionNumber)
	})
	t.Run("Doesn't get nonexistent section", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := section.NewRepository(db)

		_, err := repo.Get(context.TODO(), 9999)
		assert.Error(t, err)
	})
}

func TestRepoGetAll(t *testing.T) {
	t.Run("Gets all sections", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := section.NewRepository(db)
		expected := getTestSection()

		repo.Save(context.TODO(), expected)

		received, err := repo.GetAll(context.TODO())
		assert.NoError(t, err)

		assert.True(t, len(received) > 0)
	})
}

func TestRepoUpdate(t *testing.T) {
	t.Run("Updates a section", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := section.NewRepository(db)
		wh := getTestSection()

		id, _ := repo.Save(context.TODO(), wh)

		wh.ID = id
		wh.SectionNumber = 99

		err := repo.Update(context.TODO(), wh)
		assert.NoError(t, err)

		received, _ := repo.Get(context.TODO(), id)
		assert.Equal(t, wh.SectionNumber, received.SectionNumber)
	})
}

func TestRepoDelete(t *testing.T) {
	t.Run("Deletes a section", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := section.NewRepository(db)
		sec := getTestSection()

		id, _ := repo.Save(context.TODO(), sec)

		err := repo.Delete(context.TODO(), id)
		assert.NoError(t, err)

		_, err = repo.Get(context.TODO(), id)
		assert.Error(t, err)
	})
}

func getTestSection() domain.Section {
	return domain.Section{
		SectionNumber:      18,
		CurrentTemperature: 30,
		MinimumTemperature: 0,
		CurrentCapacity:    23,
		MinimumCapacity:    0,
		MaximumCapacity:    100,
		WarehouseID:        1,
		ProductTypeID:      1,
	}
}
