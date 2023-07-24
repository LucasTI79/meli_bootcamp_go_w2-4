package employee_test

import (
	"context"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/employee"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestRepoCreate(t *testing.T) {
	t.Run("Creates valid employee", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := employee.NewRepository(db)
		emp := getTestEmployee()

		id, err := repo.Save(context.TODO(), emp)

		assert.NoError(t, err)

		var receivedCardID string
		row := db.QueryRow(`SELECT card_number_id FROM employees WHERE id = ?;`, id)
		row.Scan(&receivedCardID)

		assert.Equal(t, emp.CardNumberID, receivedCardID)
	})
}

func TestRepoExists(t *testing.T) {
	t.Run("Finds existing employee", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := employee.NewRepository(db)
		expected := getTestEmployee()

		repo.Save(context.TODO(), expected)

		exists := repo.Exists(context.TODO(), expected.CardNumberID)

		assert.True(t, exists)
	})
}

func TestRepoGet(t *testing.T) {
	t.Run("Gets employee", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := employee.NewRepository(db)
		expected := getTestEmployee()

		id, _ := repo.Save(context.TODO(), expected)

		received, err := repo.Get(context.TODO(), id)
		assert.NoError(t, err)

		assert.Equal(t, expected.CardNumberID, received.CardNumberID)
	})
	t.Run("Doesn't get nonexistent employee", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := employee.NewRepository(db)

		_, err := repo.Get(context.TODO(), 9999)
		assert.Error(t, err)
	})
}

func TestRepoGetAll(t *testing.T) {
	t.Run("Gets all employees", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := employee.NewRepository(db)
		emp := getTestEmployee()

		repo.Save(context.TODO(), emp)

		received, err := repo.GetAll(context.TODO())
		assert.NoError(t, err)

		assert.True(t, len(received) > 0)
	})
}

func TestRepoUpdate(t *testing.T) {
	t.Run("Updates an employee", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := employee.NewRepository(db)
		emp := getTestEmployee()

		id, _ := repo.Save(context.TODO(), emp)

		emp.ID = id
		emp.FirstName = "Joao"

		err := repo.Update(context.TODO(), emp)
		assert.NoError(t, err)

		received, _ := repo.Get(context.TODO(), id)
		assert.Equal(t, emp.FirstName, received.FirstName)
	})
}

func TestRepoDelete(t *testing.T) {
	t.Run("Deletes a employee", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := employee.NewRepository(db)
		emp := getTestEmployee()

		id, _ := repo.Save(context.TODO(), emp)

		err := repo.Delete(context.TODO(), id)
		assert.NoError(t, err)

		_, err = repo.Get(context.TODO(), id)
		assert.Error(t, err)
	})
}

func TestRepoReport(t *testing.T) {
	t.Run("Gets single report", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := employee.NewRepository(db)
		emp := getTestEmployee()

		id, _ := repo.Save(context.TODO(), emp)

		_, err := repo.GetInboundReport(context.TODO(), id)
		assert.NoError(t, err)
	})
	t.Run("Gets report for every ID", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := employee.NewRepository(db)

		_, err := repo.GetAllInboundReports(context.TODO())
		assert.NoError(t, err)
	})
}

func getTestEmployee() domain.Employee {
	return domain.Employee{
		CardNumberID: "1234",
		FirstName:    "Jane",
		LastName:     "Doe",
		WarehouseID:  1,
	}
}
