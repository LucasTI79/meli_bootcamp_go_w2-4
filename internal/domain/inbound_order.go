package domain

import "time"

type InboundOrder struct {
	ID             int       `json:"id"`
	OrderDate      time.Time `json:"order_date"`
	OrderNumber    string    `json:"order_number"`
	EmployeeID     int       `json:"employee_id"`
	ProductBatchID int       `json:"product_batch_id"`
	WarehouseID    int       `json:"warehouse_id"`
}
