package db_test

import (
	"database/sql"
	"expense-tracker/internal/db"
	"expense-tracker/internal/models"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestRepo(t *testing.T) *db.ExpenseRepo {
	sqlDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	//open in memory
	_, err = sqlDB.Exec(`CREATE TABLE expenses(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	description TEXT,
	amount REAL,
	category TEXT,
	date TEXT)`)

	if err != nil {
		t.Fatal(err)
	}
	return db.NewExpenseRepo(sqlDB)
}

func TestExpenseRepo_AddAndListExpense(t *testing.T) {
	repo := setupTestRepo(t)

	expense := models.Expense{
		Description: "Coffee",
		Amount:      3.5,
		Category:    "Food",
		Date:        "2026-01-29",
	}

	if err := repo.AddExpense(expense); err != nil {
		t.Fatalf("AddExpense failed: %v", err)
	}

	list, err := repo.ListExpense()
	if err != nil {
		t.Fatalf("List Expense failed: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("expected 1, got %d", len(list))
	}

	if list[0].Description != expense.Description {
		t.Errorf("expected description %q, got %q", expense.Description, list[0].Description)
	}
}

func TestExpenseRepo_DeleteExpense(t *testing.T) {
	repo := setupTestRepo(t)

	expense := models.Expense{
		Description: "Book",
		Amount:      10,
		Category:    "Education",
		Date:        "2026-01-29",
	}

	if err := repo.AddExpense(expense); err != nil {
		t.Fatal(err)

	}
	list, err := repo.ListExpense()
	if err != nil {
		t.Fatalf("List expense failed: %v", err)
	}
	id := list[0].ID

	if err := repo.DeleteExpense(id); err != nil {
		t.Fatalf("Delete expense failed: %v", err)
	}
	list, _ = repo.ListExpense()

	if len(list) != 0 {
		t.Fatalf("expected 0 expenses after delete, got :%d", len(list))
	}
}

func TestExpenseRepo_GetCategoryTotals(t *testing.T) {
	repo := setupTestRepo(t)

	repo.AddExpense(models.Expense{Description: "Coffee", Amount: 3.5, Category: "Food", Date: "2026-01-29"})
	repo.AddExpense(models.Expense{Description: "Lunch", Amount: 7, Category: "Food", Date: "2026-01-29"})
	repo.AddExpense(models.Expense{Description: "Book", Amount: 10, Category: "Education", Date: "2026-01-29"})

	totals, err := repo.GetCategoryTotals()
	if err != nil {
		t.Fatal(err)
	}

	if totals["Food"] != 10.5 {
		t.Fatalf("expected food total of 10.5, got %v", totals["Food"])
	}

	if totals["Education"] != 10 {
		t.Fatalf("expected education total of 10, got %v", totals["Education"])
	}
}
