# DTORM - Database ORM for Go

DTORM is a versatile and lightweight ORM (Object-Relational Mapping) package for Go that supports multiple database backends including PostgreSQL, MySQL, SQLite, and Microsoft SQL Server.

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
go get github.com/markoxley/dtorm
```

## Quick Start

```go
package main

import (
    "github.com/markoxley/dtorm"
    "github.com/markoxley/dtorm/where"
)

// Define your model
type User struct {
    dtorm.Model
    Username  string `dtorm:"username"`
    Email     string `dtorm:"email"`
}

func main() {
    // Initialize database connection
    config := dtorm.Config{
        Driver:   "mysql",
        Host:     "localhost",
        Port:     3306,
        Database: "mydb",
        Username: "user",
        Password: "password",
    }

    db, err := dtorm.Connect(config)
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

DTORM provides a powerful WHERE clause builder with support for various conditions:

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

DTORM uses struct tags to define model properties:

- `dtorm:"field_name"` - Specify database field name
- `dtorm:"field_name,pk"` - Mark field as primary key
- `dtorm:"field_name,size:255"` - Set field size
- `dtorm:"field_name,nullable"` - Allow NULL values
- `dtorm:"field_name,unique"` - Enforce uniqueness
- `dtorm:"field_name,index"` - Create index on field

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
