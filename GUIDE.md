# Mud User Guide

Mud is a versatile and lightweight ORM (Object-Relational Mapping) package for Go. It supports **SQLite**, **MySQL/MariaDB** and **Microsoft SQL Server**, providing a fluent interface for database operations.

## Installation

To install Mud, use `go get`:

```bash
go get github.com/markoxley/mud
```

## Connecting to a Database

The core of Mud is the `DB` struct. You initialize a connection using `mud.New()` with a `mud.Config` object.

### Configuration

The `mud.Config` struct supports the following fields:

| Field                  | Type   | Description                                                                                                           |
| ---------------------- | ------ | --------------------------------------------------------------------------------------------------------------------- |
| `Type`                 | string | Database type. Supported values: `"sqlite"`, `"sqlite3"`, `"mysql"`, `"mariadb"`, `"sqlserver"`, `"mssql"`.           |
| `Host`                 | string | Hostname or IP (optional for SQLite).                                                                                 |
| `Database`             | string | Database name (or file path for SQLite).                                                                              |
| `User`                 | string | Username (optional for SQLite).                                                                                       |
| `Password`             | string | Password (optional for SQLite).                                                                                       |
| `Deletable`            | bool   | If `true`, `Remove` performs hard deletes (SQL DELETE). If `false` (default), soft deletes are used (DeleteDate set). |
| `DisabledTransactions` | bool   | If `true`, disables automatic transaction handling.                                                                   |

### Example Connection

```go
package main

import (
    "log"
    "github.com/markoxley/mud"
)

func main() {
    config := &mud.Config{
        Type:     "mysql",
        Host:     "localhost:3306",
        Database: "myapp",
        User:     "admin",
        Password: "securepassword",
    }

    db, err := mud.New(config)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // DB is ready to use
}
```

## Defining Models

Models in Mud are Go structs that implement the `Modeller` interface. The easiest way to do this is to embed `mud.Model`.

### Basic Model Structure

By default, Mud uses the **struct name** as the table name. You can override
it by implementing a `TableName() string` method on the model, or by calling
`SetTableName` on the embedded `mud.Model`:

```go
func (t *User) TableName() string { return "users" }
```

```go
type User struct {
    mud.Model // Embeds ID (UUID), CreateDate, LastUpdate, DeleteDate

    Username string `mud:"size:64,key:true"`
    Email    string `mud:"size:256"`
    Age      int
    Active   bool
}
```

### Model Tags (`mud:"..."`)

Use the `mud` struct tag to configure fields. Multiple options are comma-separated.

- **`key:true`**: Creates a database index on the field.
- **`identity:true`**: Marks the field as an identity column.
- **`size:N[,M]`**: Sets field size.
  - For strings/char: `N` is length.
  - For decimals: `N` is precision, `M` is scale.
- **`allowNull`**: Allows the database column to be NULL.
- **`unsigned`**: Marks numeric types as unsigned.
- **`type:X`**: Overrides the detected type. Common values: `text`, `time`.

### IDs and Timestamps

When embedding `mud.Model`:

- **ID**: A UUID string is automatically generated for new records.
- **CreateDate**: Set automatically on creation.
- **LastUpdate**: Updated automatically on save.
- **DeleteDate**: Used for soft deletes.

## CRUD Operations

### Creating and Updating

Mud uses a single `Save` method for both creating (INSERT) and updating (UPDATE) records. It determines the operation based on whether the `ID` is nil (New) or set (Update).

```go
// Create
user := &User{
    Username: "jdoe",
    Email:    "jdoe@example.com",
}
err := db.Save(user) // user.ID is now populated with a new UUID

// Update
user.Email = "john.doe@example.com"
err = db.Save(user) // Updates the existing record
```

### Reading Data

Mud provides generic helper functions for type-safe querying.

#### Fetching by ID

```go
user, err := mud.FromID[User](db, "some-uuid-string")
```

#### Fetching a Single Record

```go
// Returns the first user named 'jdoe'
user, err := mud.First[User](db, where.Equal("Username", "jdoe"))
```

#### Fetching Multiple Records

```go
// Returns all users over age 18
users, err := mud.Fetch[User](db, where.Greater("Age", 18))
```

#### Iterating (Range)

For large datasets, use `Range` to iterate over results one by one efficiently.

```go
for user := range mud.Range[User](db, where.Greater("Age", 18)) {
    fmt.Printf("User: %s\n", user.Username)
}
```

#### Counting

```go
count := db.Count(&User{}, where.Equal("Active", true))
```

#### Refreshing

To reload a model's data from the database:

```go
err := db.Refresh(user)
```

### Deleting

#### Single Record

```go
err := db.Remove(user)
```

_Note: By default, this performs a **soft delete** (sets `DeleteDate`). Set `Config.Deletable = true` to perform hard deletes._

#### Multiple Records

```go
count, err := db.RemoveMany(&User{}, where.Equal("Active", false))
```

## Query Building

Mud features a powerful, fluent query builder via the `where` and `order` packages.

### Where Clauses

Import the `where` package: `import "github.com/markoxley/mud/where"`

```go
// Basic comparison
where.Equal("Status", "active")
where.NotEqual("Status", "banned")
where.Greater("Age", 21)
where.Less("Price", 100)
where.Between("Price", 10, 100)

// String operations
where.Like("Name", "%John%")
where.Contains("Name", "John")
where.StartsWith("SKU", "PROD-")
where.EndsWith("Email", "@company.com")

// Lists
where.In("Category", []interface{}{"A", "B", "C"})
where.NotIn("Id", []interface{}{1, 2, 3})

// Null checks
where.IsNull("DeletedAt")
where.IsNotNull("Avatar")

// Chaining (AND)
where.Equal("Status", "active").AndGreater("Age", 21)

// Chaining (OR)
where.Equal("Role", "admin").OrEqual("Role", "moderator")

// Complex Grouping
// (Status = 'active' AND Age > 21) OR (Role = 'admin')
where.Equal("Status", "active").AndGreater("Age", 21).OrSub(
    where.Equal("Role", "admin"),
)
```

### Sorting

Import the `order` package: `import "github.com/markoxley/mud/order"`

```go
// Sort by Name ascending
o := order.Asc("Name")

// Sort by Age descending, then Name ascending
o := order.Desc("Age").Asc("Name")
```

## Advanced Features

### Standing Data (Seeding)

You can define default data that should exist in the table by implementing the `StandingData()` method on your model. When Mud checks the table definition (e.g., on first access), it will automatically insert this data if the table is created.

```go
func (u User) StandingData() []mud.Modeller {
    return []mud.Modeller{
        &User{Username: "admin", Role: "admin"},
        &User{Username: "guest", Role: "guest"},
    }
}
```

### Lifecycle Hooks

Mud supports hooks via interfaces:

- **`Updater`**: Implement this interface to execute logic before a `Save`.
  ```go
  func (u *User) Update(mgr mud.Manager) error {
      // logic before update/insert
      return nil
  }
  ```
- **`Restorer`**: Implement this interface to execute logic after a model is populated from the database.
  ```go
  func (u *User) Restore(mgr mud.Manager) error {
      // logic after fetch
      return nil
  }
  ```

### Transactions

Mud supports explicit transactions.

```go
tx, err := db.BeginTransaction()
if err != nil {
    // handle error
}

// Pass the transaction to operations as a variadic argument
err = db.Save(user, tx)
if err != nil {
    db.RollbackTransaction(tx)
    return
}

db.CommitTransaction(tx)
```

### Contexts and Timeouts

Use `WithContext` to bind a `context.Context` — for cancellation or a
deadline — to all operations performed through the returned clone. The clone
is a regular `*DB` that shares the parent's connection pool and caches, so it
can be passed to any function (including the generic helpers) that takes a
`*DB`. Do not call `Close` on a clone — it shares the parent's pool.

```go
// Per-request context (e.g. an HTTP handler)
users, err := mud.Fetch[User](db.WithContext(r.Context()))

// Query timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
rows, err := db.WithContext(ctx).Fetch(&User{})

// Reusable handle
cdb := db.WithContext(ctx)
err = cdb.Save(&user)
```

The deadline applies to the whole operation, including any transaction Mud
opens for it (`BeginTx`). A `nil` context is treated as
`context.Background()`, and the parent `*DB` is unaffected by the clone's
context.

### Raw SQL

For complex operations not supported by the query builder. All raw methods
accept a slice of arguments which are bound safely as query parameters —
always prefer them over interpolating values into the SQL string.

```go
// Execute non-query (UPDATE, DELETE, etc.)
err := db.RawExecute("UPDATE users SET legacy = ? WHERE created_at < ?", []interface{}{1, "2020-01-01"})

// Query for results (returns []map[string]interface{})
results, err := db.RawSelect("SELECT count(*) as cnt, status FROM users GROUP BY status", nil)

// Scalar queries also take bind arguments
val, ok := db.RawScalar("SELECT count(*) FROM users WHERE status = ?", []interface{}{"active"})
```

_Note: SQL Server uses `@p1`, `@p2`, ... placeholders instead of `?`._

### Schema Management

Mud automatically attempts to ensure the database schema matches your models. When you perform an operation on a model, Mud checks if the table exists. If not, it generates a `CREATE TABLE` statement based on the struct fields and tags, and creates the table (and indexes).
