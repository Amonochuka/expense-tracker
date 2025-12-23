package db

import (
	"database/sql"
	"expense-tracker/internal/models"
	"fmt"
)

type ExpenseRepo struct {
	db *sql.DB
}

func NewExpenseRepo(db *sql.DB) *ExpenseRepo {
	return &ExpenseRepo{db: db}
}

func (r *ExpenseRepo) AddExpense(e models.Expense) error {
	_, err := r.db.Exec(`INSERT INTO expenses(description, amount, category, date)
		 VALUES(?, ?, ?, ?)`, e.Description, e.Amount, e.Category, e.Date)
	if err != nil {
		fmt.Println("DB Exec error:", err)
		return err
	}
	return nil
}

func (r *ExpenseRepo) ListExpense() ([]models.Expense, error) {
	rows, err := r.db.Query(`SELECT id, description, amount, category, date FROM expenses`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var expenses []models.Expense
	for rows.Next() {
		var e models.Expense
		if err := rows.Scan(&e.ID, &e.Description, &e.Amount, &e.Category, &e.Date); err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}
	return expenses, nil
}

func (r *ExpenseRepo) DeleteExpense(id int64) error {
	res, err := r.db.Exec(`DELETE FROM expenses WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *ExpenseRepo) FilterExpensesByCategory(category string) ([]models.Expense, error) {
	rows, err := r.db.Query(`SELECT id, description, amount, category, date FROM expenses 
	WHERE category = ? 
	ORDER BY DATE`, category)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var expenses []models.Expense
	for rows.Next() {
		var e models.Expense
		if err := rows.Scan(&e.ID, &e.Description, &e.Amount, &e.Category, &e.Date); err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}
	return expenses, nil
}

func (r *ExpenseRepo) FilterExpensesByDate(start, end string) ([]models.Expense, error) {
	//validate time
	rows, err := r.db.Query(`SELECT id, description, amount, category, date FROM expenses 
              WHERE date BETWEEN ? AND ? ORDER BY date`, start, end)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var e models.Expense
		if err := rows.Scan(&e.ID, &e.Description, &e.Amount, &e.Category, &e.Date); err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}
	return expenses, nil
}

/*
Why this is different
You do not have a struct that represents:
{ category, sum }
A map entry like totals[category]:
Does not exist until you assign it, it's not addressable, So you must:
Scan into variables, Then assign into the map
*/
func (r *ExpenseRepo) GetCategoryTotals() (map[string]float64, error) {
	rows, err := r.db.Query(`SELECT category, SUM(amount) FROM expenses GROUP BY category`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	totals := make(map[string]float64)
	for rows.Next() {
		var category string
		var sum float64
		if err := rows.Scan(&category, &sum); err != nil {
			return nil, err
		}
		totals[category] = sum
	}
	return totals, nil
}
