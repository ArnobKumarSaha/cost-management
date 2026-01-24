package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// importFile reads the file and saves to DB
func importFile(filename string, defaultYear int, collection *mongo.Collection) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentDate time.Time
	var expenses []Expense

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) >= 2 {
			dayStr := parts[0]
			if day, err := strconv.Atoi(dayStr); err == nil {
				month, err := parseMonth(parts[1])
				if err == nil {
					// New date found → save previous day first (if any)
					if !currentDate.IsZero() && len(expenses) > 0 {
						saveDailyExpenses(collection, currentDate, expenses)
					}
					// Reset for new day
					expenses = []Expense{}
					currentDate = time.Date(defaultYear, month, day, 0, 0, 0, 0, time.UTC)
					continue
				}
			}
		}

		// It's an expense line
		if currentDate.IsZero() {
			continue
		}

		exp, err := parseExpense(line)
		if err != nil {
			fmt.Printf("Warning: skipping line '%s': %v\n", line, err)
			continue
		}
		expenses = append(expenses, exp)
	}

	// Don't forget the very last day
	if !currentDate.IsZero() && len(expenses) > 0 {
		saveDailyExpenses(collection, currentDate, expenses)
	}

	return scanner.Err()
}

// parseMonth converts month abbr to time.Month
func parseMonth(abbr string) (time.Month, error) {
	abbr = strings.ToLower(abbr)
	months := map[string]time.Month{
		"jan": time.January,
		"feb": time.February,
		"mar": time.March,
		"apr": time.April,
		"may": time.May,
		"jun": time.June,
		"jul": time.July,
		"aug": time.August,
		"sep": time.September,
		"oct": time.October,
		"nov": time.November,
		"dec": time.December,
	}
	m, ok := months[abbr]
	if !ok {
		return 0, fmt.Errorf("invalid month: %s", abbr)
	}
	return m, nil
}

// parseExpense parses "Item with spaces amount cat"
func parseExpense(line string) (Expense, error) {
	parts := strings.Split(line, " ")
	if len(parts) < 3 {
		return Expense{}, fmt.Errorf("invalid expense format")
	}
	cat := parts[len(parts)-1]
	amountStr := parts[len(parts)-2]
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return Expense{}, err
	}
	item := strings.Join(parts[:len(parts)-2], " ")

	fullCat, ok := CategoryMap[cat]
	if !ok {
		fullCat = cat // Use as is if not mapped
	}

	return Expense{Item: item, Amount: amount, Category: fullCat}, nil
}

// saveDailyExpenses upserts the daily expenses (overwrite if exists)
func saveDailyExpenses(collection *mongo.Collection, date time.Time, expenses []Expense) {
	ctx := context.TODO()
	filter := bson.M{"date": date}
	update := bson.M{"$set": bson.M{"expenses": expenses}}
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		fmt.Printf("Error saving date %s: %v\n", date.Format("2006-01-02"), err)
	} else {
		fmt.Printf("Saved/Updated expenses for %s\n", date.Format("2006-01-02"))
	}
}
