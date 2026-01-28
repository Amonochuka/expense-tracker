package models

type Expense struct {
	ID          int64
	Description string
	Amount      float64
	Category    string
	Date        string // YYYY-MM-DD
}
