package models

type Expense struct {
	ID          int
	Description string
	Amount      float64
	Category    string
	Date        string // YYYY-MM-DD
}
