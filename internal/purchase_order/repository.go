package purchaseorder

import (
	"context"
	"database/sql"
	"strings"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, i domain.PurchaseOrder) (int, error)
	Exists(ctx context.Context, orderNumber string) bool
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Exists(ctx context.Context, orderNumber string) bool {
	query := "SELECT order_number FROM purchase_orders WHERE order_number=?;"
	row := r.db.QueryRow(query, orderNumber)
	err := row.Scan(&orderNumber)
	return err == nil
}

func (r *repository) Create(ctx context.Context, i domain.PurchaseOrder) (int, error) {
	queryPurchaseOrders := "INSERT INTO purchase_orders(order_number,order_date,tracking_code,buyer_id,order_status_id,product_record_id) SELECT ?,?,?,?,?,? FROM product_records pr WHERE pr.id = ?"
	res, err := r.db.Exec(queryPurchaseOrders, i.OrderNumber, i.OrderDate, i.TrackingCode, i.BuyerID, i.OrderStatusID, i.ProductRecordID, i.ProductRecordID)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1452") {
			return 0, ErrFKNotFound
		}
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	queryOrderDetails := "INSERT INTO order_details(clean_liness_status,quantity,temperature,product_record_id,purchase_order_id) VALUES (?,?,?,?,?)"
	res, err = r.db.Exec(queryOrderDetails, "good", 1, 32, i.ProductRecordID, id)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1452") {
			return 0, ErrProductRecordIDNotFound
		}
		return 0, err
	}
	return int(id), nil
}
