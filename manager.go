package dtorm

import (
	"fmt"
)

type Manager interface {
	ConnectionString(cfg *Config) (string, error)
	// WhereString(c *Criteria) (string, error)
	// OrderString(c *Criteria) (string, error)
	LimitString(c *Criteria) string
	OffsetString(c *Criteria) string
	IdentityString(f string) string
	TableDefinition(m Modeller) ([]string, bool)
	InsertCommand(m Modeller) string
	UpdateCommand(m Modeller) string
	DeleteCommand(m Modeller) string
	RefreshCommand(m Modeller) string
	BuildQuery(where string, order string, limit string, offset string) string
	TableTest(m Modeller) ([]field, string, bool)
	TableExistsQuery(dbName, name string) string
	MassDelete(m Modeller, c *Criteria) string
	MassDisable(m Modeller, c *Criteria) string
	Operators() []string
	TableCreate() string
	IndexCreate() string
}

func GetManager(config *Config) (Manager, error) {
	switch config.Type {
	// case "sqlite":
	// 	return &SqliteManager{}, nil
	case "mysql":
		return &MySQLManager{}, nil
	// case "sqlserver":
	// 	return &MSSQLManager{}, nil
	default:
		return nil, fmt.Errorf("invalid database type: %s", config.Type)
	}
}
