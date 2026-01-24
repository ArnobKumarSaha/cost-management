package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// generateReport generates and prints a cost sheet
func generateReport(period string, year int, month int, week int, collection *mongo.Collection) error {
	ctx := context.TODO()
	var start, end time.Time
	var title string

	switch period {
	case "yearly":
		start = time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		end = time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
		title = fmt.Sprintf("Yearly Cost Sheet for %d", year)
	case "monthly":
		if month < 1 || month > 12 {
			return fmt.Errorf("invalid month: %d", month)
		}
		start = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		end = start.AddDate(0, 1, 0)
		title = fmt.Sprintf("Monthly Cost Sheet for %s %d", time.Month(month), year)
	case "weekly":
		if week < 1 || week > 53 {
			return fmt.Errorf("invalid week: %d", week)
		}
		// Approximate: start from Jan 1, add (week-1)*7 days
		jan1 := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		start = jan1.AddDate(0, 0, (week-1)*7)
		end = start.AddDate(0, 0, 7)
		title = fmt.Sprintf("Weekly Cost Sheet for Week %d of %d (%s to %s)", week, year, start.Format("2006-01-02"), end.AddDate(0, 0, -1).Format("2006-01-02"))
	default:
		return fmt.Errorf("invalid period: %s", period)
	}

	// Query documents in date range
	filter := bson.M{
		"date": bson.M{
			"$gte": start,
			"$lt":  end,
		},
	}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	categoryTotals := make(map[string]float64)
	total := 0.0

	for cursor.Next(ctx) {
		var daily DailyExpenses
		if err := cursor.Decode(&daily); err != nil {
			fmt.Printf("Warning: decode error: %v\n", err)
			continue
		}
		for _, exp := range daily.Expenses {
			categoryTotals[exp.Category] += exp.Amount
			total += exp.Amount
		}
	}

	if err := cursor.Err(); err != nil {
		return err
	}

	// Print the report
	fmt.Println(title)
	fmt.Println("-----------------------------")
	for cat, amt := range categoryTotals {
		fmt.Printf("%s: %.2f\n", cat, amt)
	}
	fmt.Println("-----------------------------")
	fmt.Printf("Total: %.2f\n", total)

	return nil
}
