# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

A minimal Go CLI tool for personal expense tracking backed by MongoDB Atlas. It has two modes: importing expenses from a plain-text file and generating cost reports.

## Setup

Create a `.env` file in the project root:

```
MONGODB_URI="mongodb+srv://<user>:<pass>@arnob.0b4vj7e.mongodb.net/?retryWrites=true&w=majority"
```

## Commands

```bash
# Import expenses from a text file
go run *.go --mode=import --file=january.txt

# Generate a monthly report (month=1 for January)
go run *.go --mode=report --period=monthly --month=1

# Generate a weekly report
go run *.go --mode=report --period=weekly --week=3

# Generate a yearly report
go run *.go --mode=report --period=yearly --year=2026
```

## Architecture

All code lives in a single `main` package with three files:

- [main.go](main.go) — entry point, flag parsing, MongoDB connection, mode dispatch
- [import.go](import.go) — parses the expense text file and upserts records into MongoDB
- [report.go](report.go) — queries MongoDB by date range and aggregates by category

**Data model** (MongoDB, database `expenses`, collection `daily`):
- Each document is a `DailyExpenses`: a date + a list of `Expense` items (item name, amount, category)
- Import upserts by date (overwrite if the day already exists)

**Input file format** (`january.txt`):
```
18 jan
Rickshaw 30 Tr
Cha 20 Sn
```
Lines starting with `<day> <month-abbr>` start a new day. Expense lines are `<item name> <amount> <category-abbr>`. Category abbreviations are mapped in `CategoryMap` in [main.go](main.go#L30).
