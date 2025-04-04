# mud - Database ORM for Go

_Note: This is a work in progress and is not yet ready for production use._

mud is a versatile and lightweight ORM (Object-Relational Mapping) package for Go that supports multiple database backends including PostgreSQL, MySQL, SQLite, and Microsoft SQL Server.

## Features

- Support for multiple database backends
- Fluent query builder interface
- Advanced WHERE clause construction
- Automatic model mapping
- Transaction support
- Field validation and type safety
- UUID support for primary keys

## Installation

```bash
go get github.com/markoxley/mud
```

## Quick Start

```go
package main

import (
    "github.com/markoxley/mud"
    "github.com/markoxley/mud/where"
)

// Define your model
type User struct {
    mud.Model
    Username  string `mud:"size:64,key:true"`
    Email     string `mud:"size:256"`
}

func main() {
    // Initialize database connection
    config := mud.Config{
        Driver:   "mysql",
        Host:     "localhost",
        Port:     3306,
        Database: "mydb",
        Username: "user",
        Password: "password",
    }

    db, err := mud.Connect(config)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    // Create a new user
    user := &User{
        Username: "johndoe",
        Email:    "john@example.com",
    }

    err = db.Save(user)
    if err != nil {
        panic(err)
    }

    // Query users with WHERE clause
    var users []User
    where := where.Equal("username", "johndoe")
    err = db.Find(&users, where)
    if err != nil {
        panic(err)
    }

    // Update user, we reuse the Save method
    user.Email = "john.doe@example.com"
    err = db.Save(user)
    if err != nil {
        panic(err)
    }

    // Delete user
    err = db.Remove(user)
    if err != nil {
        panic(err)
    }

    // Convenience functions with generics
    user, err = mud.First[User](db, where.Equal("username", "johndoe"))
    if err != nil {
        panic(err)
    }

    user, err = mud.FromID[User](db, *user.ID)
    if err != nil {
        panic(err)
    }
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
where.NotIsNull("field")

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
- `mud:"size:255"` - Set field size
- `mud:"allowNull"` - Allow NULL values

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
