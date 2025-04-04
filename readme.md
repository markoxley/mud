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
    Username  string `mud:"username"`
    Email     string `mud:"email"`
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

    err = db.Create(user)
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

## Model Tags

mud uses struct tags to define model properties:

- `mud:"field_name"` - Specify database field name
- `mud:"field_name,pk"` - Mark field as primary key
- `mud:"field_name,size:255"` - Set field size
- `mud:"field_name,nullable"` - Allow NULL values
- `mud:"field_name,unique"` - Enforce uniqueness
- `mud:"field_name,index"` - Create index on field

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
