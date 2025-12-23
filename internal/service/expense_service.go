package service

import (
	"errors"
	"expense-tracker/internal/db"
	"expense-tracker/internal/models"
	"time"
)

type ExpenseService struct {
	repo *db.ExpenseRepo
}

func NewExpenseService(repo *db.ExpenseRepo) *ExpenseService {
	return &ExpenseService{repo: repo}
}

func (s *ExpenseService) AddExpense(exp models.Expense) error {
	if exp.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if exp.Description == "" {
		return errors.New("description cannnot be empty")

	}
	if exp.Category == "" {
		return errors.New("category cannot be empty")
	}

	if _, err := time.Parse("2006-01-02", exp.Date); err != nil {
		return errors.New("invalid date format")
	}
	return s.repo.AddExpense(exp)
}

func (s *ExpenseService) ListExpense() ([]models.Expense, error) {
	return s.repo.ListExpense()
}

func (s *ExpenseService) DeleteExpense(id int64) error {
	return s.repo.DeleteExpense(id)
}

func (s *ExpenseService) FilterExpensesByCategory(category string) ([]models.Expense, error) {
	return s.repo.FilterExpensesByCategory(category)
}

func (s *ExpenseService) FilterExpensesByDate(start, end string) ([]models.Expense, error) {
	return s.repo.FilterExpensesByDate(start, end)
}

func (s *ExpenseService) GetCategoryTotals() (map[string]float64, error) {
	return s.repo.GetCategoryTotals()
}
