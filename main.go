package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Expense represents a single expense entry
type Expense struct {
	Item     string  `bson:"item"`
	Amount   float64 `bson:"amount"`
	Category string  `bson:"category"`
}

// DailyExpenses represents expenses for a single day
type DailyExpenses struct {
	Date     time.Time `bson:"date"`
	Expenses []Expense `bson:"expenses"`
}

// CategoryMap maps abbreviations to full names
var CategoryMap = map[string]string{
	"Tr": "Transport",
	"Bz": "Groceries",
	"Sn": "Snacks",
	"Ot": "OneTime",
	"Fd": "Food",
	"Ut": "Utility",
	"Ex": "Extras",
	"To": "Tour",
}

func main() {
	// Command-line flags
	mode := flag.String("mode", "import", "Mode: 'import' to read file and save to DB, 'report' to generate cost sheet")
	file := flag.String("file", "", "Input file for import mode")
	period := flag.String("period", "monthly", "Period for report: 'weekly', 'monthly', 'yearly'")
	year := flag.Int("year", time.Now().Year(), "Year for dates (default current year)")
	month := flag.Int("month", 0, "Month for monthly report (1-12)")
	week := flag.Int("week", 0, "Week number for weekly report (1-53)")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: Could not load .env file")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		fmt.Println("Error: MONGODB_URI not found in .env")
		return
	}

	// MongoDB connection (replace with your Atlas URI, use env var in production)
	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		fmt.Printf("Error connecting to MongoDB: %v\n", err)
		return
	}
	defer client.Disconnect(context.TODO())
	fmt.Println("Connected to MongoDB")

	db := client.Database("expenses")
	collection := db.Collection("daily")

	switch *mode {
	case "import":
		if *file == "" {
			fmt.Println("File is required for import mode")
			return
		}
		err := importFile(*file, *year, collection)
		if err != nil {
			fmt.Printf("Error importing file: %v\n", err)
		}
	case "report":
		err := generateReport(*period, *year, *month, *week, collection)
		if err != nil {
			fmt.Printf("Error generating report: %v\n", err)
		}
	default:
		fmt.Println("Invalid mode. Use 'import' or 'report'")
	}
}
