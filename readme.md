# mud — a lightweight ORM for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/markoxley/mud.svg)](https://pkg.go.dev/github.com/markoxley/mud)
[![Go Report Card](https://goreportcard.com/badge/github.com/markoxley/mud)](https://goreportcard.com/report/github.com/markoxley/mud)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`mud` is a lightweight ORM for Go supporting **SQLite**, **MySQL/MariaDB** and **Microsoft SQL Server**. Models are plain structs — tables, columns and indexes are generated automatically from `mud` struct tags, and all queries are fully parameterized for the target dialect.

> **Status:** pre-1.0 — the API may still change. See [CHANGELOG.md](CHANGELOG.md) for recent changes and breaking items.

## Features

- Multi-dialect SQL generation — correct column types and placeholders (`?` vs `@p1`) per engine
- Fully parameterized queries — no string-interpolated SQL
- Fluent `where` and `order` clause builders
- Struct-tag–driven schema: automatic `CREATE TABLE` and index creation
- One `Save` for INSERT/UPDATE; soft deletes by default, hard deletes via `Config.Deletable`
- Transactions — implicit per operation on MySQL/SQL Server, explicit `BeginTransaction` anywhere
- `context.Context` via `db.WithContext(ctx)` — cancellation and query timeouts
- Generic helpers (`mud.Fetch[T]`, `mud.First[T]`, `mud.FromID[T]`) and lazy `mud.Range[T]` iteration
- Lifecycle hooks (`Updater`/`Restorer`) and standing-data seeding
- Raw SQL escape hatches: `RawExecute`, `RawSelect`, `RawScalar`

## Requirements

Go 1.24 or later.

## Installation

```bash
go get github.com/markoxley/mud
```

## Quick Start

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/markoxley/mud"
	"github.com/markoxley/mud/order"
	"github.com/markoxley/mud/where"
)

// Models embed mud.Model — giving you ID (UUID), CreateDate,
// LastUpdate and DeleteDate for free.
type User struct {
	mud.Model
	Username string `mud:"size:64,key:true"`
	Email    string `mud:"size:256"`
	Age      int    `mud:""`
}

func main() {
	// SQLite needs no host or credentials — just a file path.
	db, err := mud.New(&mud.Config{Type: "sqlite", Database: "app.db"})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Save inserts new records and updates existing ones.
	u := &User{Username: "jdoe", Email: "jdoe@example.com", Age: 33}
	if err := db.Save(u); err != nil {
		log.Fatal(err)
	}

	// Type-safe generic reads.
	adults, err := mud.Fetch[User](db, where.Greater("Age", 18))
	if err != nil {
		log.Fatal(err)
	}
	first, err := mud.First[User](db, where.Equal("Username", "jdoe"))
	if err != nil {
		log.Fatal(err)
	}
	byID, err := mud.FromID[User](db, *first.GetID())
	if err != nil {
		log.Fatal(err)
	}

	// Combine WHERE, ORDER BY and LIMIT through Criteria.
	sorted, err := mud.Fetch[User](db, &mud.Criteria{
		Where: where.Contains("Email", "@example.com"),
		Order: order.Desc("Age"),
		Limit: 10,
	})

	// Lazy iteration over large result sets.
	for usr := range mud.Range[User](db, where.Greater("Age", 65)) {
		println(usr.Username)
	}

	// Cancellation and timeouts via context.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := db.WithContext(ctx).Fetch(User{}); err != nil {
		log.Fatal(err)
	}

	// Soft delete by default; hard delete with Config.Deletable = true.
	if err := db.Remove(u); err != nil {
		log.Fatal(err)
	}

	_, _, _ = adults, byID, sorted
}
```

## WHERE Clause Builder

mud provides a powerful WHERE clause builder with support for various conditions:

```go
// Basic conditions
where.Equal("field", value)
where.NotEqual("field", value)
where.Greater("field", value)
where.Less("field", value)

// String operations
where.Contains("field", "substring")
where.StartsWith("field", "prefix")
where.EndsWith("field", "suffix")

// Null checks
where.IsNull("field")
where.IsNotNull("field")

// Range operations
where.Between("field", value1, value2)
where.In("field", []interface{}{value1, value2})

// Combining conditions
where.Equal("field1", value1).AndEqual("field2", value2)
where.Equal("field1", value1).OrEqual("field2", value2)
```

## Order clause builder

mud provides a powerful ORDER BY clause builder with support for various conditions:

```go
// Basic conditions
order.Asc("field")
order.Desc("field")

// Combining conditions
order.Asc("field1").Desc("field2")
```

## Model Tags

mud uses struct tags to define model properties:

- `mud:""` - Specify the field is to be included in the database
- `mud:"key:true"` - Create an index on field
- `mud:"size:255"` - Set field size (`size:10,2` for decimal precision/scale)
- `mud:"allowNull"` - Allow NULL values
- `mud:"unsigned"` - Mark numeric fields unsigned (widened on SQL Server)
- `mud:"identity:true"` - Mark a field as an identity column
- `mud:"type:X"` - Override the detected column type (e.g. `text`, `decimal`, `time`)

## Contexts and Timeouts

`WithContext` returns a `*DB` clone bound to a context — every operation
through it honours cancellation and deadlines, and it can be passed to any
function that takes a `*DB` (including the generic helpers):

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

users, err := mud.Fetch[User](db.WithContext(ctx))
```

See [GUIDE.md](GUIDE.md#contexts-and-timeouts) for details.

## Documentation

- [User Guide](GUIDE.md) — models, struct tags, query building, transactions, contexts, raw SQL
- [API Reference](https://pkg.go.dev/github.com/markoxley/mud) on pkg.go.dev
- [Changelog](CHANGELOG.md)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
