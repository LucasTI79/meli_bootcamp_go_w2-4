package domain

type Buyer struct {
	ID           int    `json:"id"`
	CardNumberID string `json:"card_number_id"`
	FirstName    string `binding:"required" json:"first_name"`
	LastName     string `binding:"required" json:"last_name"`
}

type BuyerCreate struct {
	ID           int    `json:"id"`
	CardNumberID string `binding:"required" json:"card_number_id"`
	FirstName    string `binding:"required" json:"first_name"`
	LastName     string `binding:"required" json:"last_name"`
}
