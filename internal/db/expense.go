package db

import (
	"expense-tracker/internal/models"
	"fmt"
	"time"
)

// add new expense
func AddExpense(e models.Expense) error {
	fmt.Printf("DEBUG: Adding expense: %+v\n", e) // debug

	query := `INSERT INTO expenses(description, amount, category, date) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, e.Description, e.Amount, e.Category, e.Date)
	if err != nil {
		fmt.Println("DEBUG: DB.Exec error:", err)
		return err
	}

	affected, _ := res.RowsAffected()
	fmt.Println("DEBUG: Rows affected:", affected)
	return nil
}

// list expense
func ListExpenses() ([]models.Expense, error) {
	query := `SELECT id, description, amount, category, date FROM expenses ORDER BY date`
	rows, err := DB.Query(query)
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

// delete expense by ID
func DeleteExpenses(id int) error {
	query := `DELETE FROM expenses WHERE id = ?`
	res, err := DB.Exec(query, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("no expense found with ID %d", id)
	}
	return nil
}

// // FilterExpensesByCategory returns expenses in a category
func FilterExpensesByCategory(category string) ([]models.Expense, error) {
	query := `SELECT id, description, amount, category FROM expenses WHERE category = ? 
	ORDER BY DATE`
	rows, err := DB.Query(query, category)
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

// filter by date
func FilterExpensesByDate(start, end string) ([]models.Expense, error) {
	//validate time
	if _, err := time.Parse("2006-01-02", start); err != nil {
		return nil, fmt.Errorf("invalid start date format :%v", err)
	}

	if _, err := time.Parse("2006-01-02", end); err != nil {
		return nil, fmt.Errorf("invalid end date format : %v", err)
	}

	query := `SELECT id, description, amount, category, date FROM expenses 
              WHERE date BETWEEN ? AND ? ORDER BY date`
	rows, err := DB.Query(query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var e models.Expense
		if err := rows.Scan(&e.ID, &e.Amount, &e.Category, &e.Date, &e.Description); err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}
	return expenses, nil
}

// GetCategoryTotals returns sum of amounts per category
func GetCategoryTotals() (map[string]float64, error) {
	query := `SELECT category, SUM(amount) FROM expenses GROUP BY category`
	rows, err := DB.Query(query)
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
