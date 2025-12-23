package main

import (
	"bufio"
	"expense-tracker/internal/db"
	"expense-tracker/internal/models"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	db.InitDB("expenses.db")

	expenseRepo := db.NewExpenseRepo(db.DB)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to Expense tracker CLI")

	for {
		fmt.Print("\nCommands: add | list | filter-category | filter-date | delete | totals | exit\n> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		switch input {
		case "add":
			addExpense(scanner, expenseRepo)
		case "list":
			listExpenses(expenseRepo)
		case "filter-category":
			filterExpensesByCategory(scanner, expenseRepo)
		case "filter-date":
			filterExpensesByDate(scanner, expenseRepo)
		case "delete":
			deleteExpenses(scanner, expenseRepo)
		case "totals":
			showTotals(expenseRepo)
		case "exit":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Unkown command")
		}
	}
}

func addExpense(scanner *bufio.Scanner, repo *db.ExpenseRepo) {
	fmt.Print("Description :")
	scanner.Scan()
	desc := scanner.Text()

	fmt.Print("Amount :")
	scanner.Scan()
	amtstr := scanner.Text()
	amount, err := strconv.ParseFloat(amtstr, 64)
	if err != nil {
		fmt.Println("Invalid Amount")
		return
	}

	if amount < 0 {
		fmt.Println("Amount cannot be negative!")
		return
	}

	fmt.Print("Category: ")
	scanner.Scan()
	category := scanner.Text()

	fmt.Print("Date (YYYY-MM-DD, leave empty for today): ")
	scanner.Scan()
	date := scanner.Text()
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	exp := models.Expense{
		Description: desc,
		Amount:      amount,
		Category:    category,
		Date:        date,
	}

	err = repo.AddExpense(exp)
	if err != nil {
		fmt.Println("Error adding expense :", err)
	} else {
		fmt.Println("Expense added successfully")
	}
}

func listExpenses(repo *db.ExpenseRepo) {
	expenses, err := repo.ListExpense()
	if err != nil {
		fmt.Println("Error listing expenses:", err)
		return
	}

	fmt.Println("\nID | Description | Amount | Category | Date")
	fmt.Println("---------------------------------------------")

	for _, e := range expenses {
		fmt.Printf("%d | %s | %.2f | %s | %s\n", e.ID, e.Description, e.Amount, e.Category, e.Date)
	}
}

func deleteExpenses(scanner *bufio.Scanner, repo *db.ExpenseRepo) {
	fmt.Print("Enter ID to delete:")
	scanner.Scan()
	idstr := scanner.Text()
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		fmt.Println("Invalid ID")
		return
	}

	err = repo.DeleteExpense(id)
	if err != nil {
		fmt.Println("Error deleting expense:", err)
	} else {
		fmt.Println("Expense deleted successfully")
	}
}

func filterExpensesByCategory(scanner *bufio.Scanner, repo *db.ExpenseRepo) {
	fmt.Print("Enter category: ")
	scanner.Scan()
	category := scanner.Text()

	expenses, err := repo.FilterExpensesByCategory(category)
	if err != nil {
		fmt.Println("Error filtering by that category !:", err)
		return
	}

	fmt.Println("\nID | Description | Amount | Category | Date")
	fmt.Println("---------------------------------------------")

	for _, e := range expenses {
		fmt.Printf("%d | %s | %.2f | %s | %s\n", e.ID, e.Description, e.Amount, e.Category, e.Date)
	}
}

func filterExpensesByDate(scanner *bufio.Scanner, repo *db.ExpenseRepo) {
	fmt.Print("Enter start date: ")
	scanner.Scan()
	start := scanner.Text()

	fmt.Print("Enter end date: ")
	scanner.Scan()
	end := scanner.Text()

	expenses, err := repo.FilterExpensesByDate(start, end)
	if err != nil {
		fmt.Println("Error filtering!:", err)
		return
	}

	fmt.Println("\nID | Description | Amount | Category | Date")
	fmt.Println("---------------------------------------------")

	for _, e := range expenses {
		fmt.Printf("%d | %s | %.2f | %s | %s\n", e.ID, e.Description, e.Amount, e.Category, e.Date)
	}
}

func showTotals(repo *db.ExpenseRepo) {
	totals, err := repo.GetCategoryTotals()
	if err != nil {
		fmt.Println("Error calculating totals:", err)
		return
	}

	fmt.Println("\nCategory | sum")
	fmt.Println("----------------")

	for cat, sum := range totals {
		fmt.Printf("%s | %.2f\n", cat, sum)
	}
}
