package domain

// Buyer represents a buyer
type Buyer struct {
	ID           int    `json:"id"`
	CardNumberID string `json:"card_number_id"`
	FirstName    string `binding:"required" json:"first_name"`
	LastName     string `binding:"required" json:"last_name"`
}

// BuyerCreate represents the data for creating a buyer
type BuyerCreate struct {
	ID           int    `json:"id"`
	CardNumberID string `binding:"required" json:"card_number_id"`
	FirstName    string `binding:"required" json:"first_name"`
	LastName     string `binding:"required" json:"last_name"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Message string `json:"message"`
}
