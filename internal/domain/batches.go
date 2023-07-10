package domain

import "time"

type Batches struct {
	ID                 int       `json:"id"`
	BatchNumber        int       `json:"batch_number"`
	CurrentQuantity    int       `json:"current_quantity"`
	CurrentTemperature int       `json:"current_temperature"`
	DueDate            time.Time `json:"due_date"`
	InitialQuantity    int       `json:"initial_quantity"`
	ManufacturingDate  time.Time `json:"manufacturing_date"`
	ManufacturingHour  int       `json:"manufacturing_hour"`
	MinimumTemperature int       `json:"minimum_temperature"`
	ProductID          int       `json:"product_id"`
	SectionID          int       `json:"section_id"`
}
