package transaction

import "time"

//Info represents a single transaction
type Info struct {
	TransactionID   int     `db:"transaction_id" json:"transaction_id"`
	CategoryID      int     `db:"category_id" json:"category_id"`
	UserID          string  `db:"user_id" json:"user_id"`
	Amount          float64 `db:"amount" json:"amount"`
	Note            string  `db:"note" json:"note"`
	TransactionDate int     `db:"transaction_date" json:"transaction_date"`
}

//NewTransaction contains information needed to create a new transaction
type NewTransaction struct {
	Amount float64   `json:"amount"`
	Note   string    `json:"note"`
	Date   time.Time `json:"transaction_date"`
}

//UpdateTransaction contains information to update a transaction
type UpdateTransaction struct {
	Amount *float64 `db:"amount" json:"amount"`
	Note   *string  `db:"note" json:"note"`
}
