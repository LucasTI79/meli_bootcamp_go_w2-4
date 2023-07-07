package inboundorder

import (
	"context"
	"database/sql"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
)

// Repository encapsulates the storage of a employee.
type Repository interface {
	Save(ctx context.Context, i domain.InboundOrder) (int, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}
func (r *repository) Save(ctx context.Context, i domain.InboundOrder) (int, error) {
	query := "INSERT INTO inbound_orders(order_date,order_number,employee_id,product_batch_id,warehouse_id) VALUES (?,?,?,?,?)"

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(&i.OrderDate, &i.OrderNumber, &i.EmployeeID, &i.ProductBatchID, &i.WarehouseID)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
