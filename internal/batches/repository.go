package batches

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, b domain.Batches) (domain.Batches, error)
	Exists(ctx context.Context, batchNumber int) bool
	Save(ctx context.Context, s domain.Batches) (int, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, b domain.Batches) (domain.Batches, error) {
	query := "INSERT INTO batches (batch_number, current_quantity, current_temperature, due_date, initial_quantity, manufacturing_date, manufacturing_hour, minimum_temperature, product_id, section_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"

	stmtIns, err := r.db.Prepare(query)
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()

	result, err := stmtIns.Exec(&b.BatchNumber, &b.CurrentQuantity, &b.CurrentTemperature, &b.DueDate, &b.InitialQuantity, &b.ManufacturingDate, &b.ManufacturingHour, &b.MinimumTemperature, &b.ProductID, &b.SectionID)
	if err != nil {
		return domain.Batches{}, err
	}

	_, err = result.LastInsertId()
	if err != nil {
		return domain.Batches{}, err
	}

	return b, nil
}

func (r *repository) Exists(ctx context.Context, batchNumber int) bool {
	query := "SELECT batch_number FROM product_batches WHERE batch_number=?;"
	row := r.db.QueryRow(query, batchNumber)
	err := row.Scan(&batchNumber)
	return err == nil
}

func (r *repository) Save(ctx context.Context, s domain.Batches) (int, error) {
	query := "INSERT INTO product_batches (batch_number, current_quantity, current_temperature, due_date, initial_quantity, manufacturing_date, manufacturing_hour, minimum_temperature, product_id, section_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}

	res, err := stmt.Exec(&s.BatchNumber, &s.CurrentQuantity, &s.CurrentTemperature, &s.DueDate, &s.InitialQuantity, &s.ManufacturingDate, &s.ManufacturingHour, &s.MinimumTemperature, &s.ProductID, &s.SectionID)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}

	return int(id), nil
}
