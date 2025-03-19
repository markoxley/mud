package dtorm

import (
	"fmt"
	"time"

	"github.com/markoxley/dtorm/utils"
)

type MySQLManager struct {
	// Needs to implement Manager
}

func (m *MySQLManager) ConnectionString(cfg *Config) (string, error) {
	if cfg == nil {
		return "", fmt.Errorf("no config provided")

	}
	if cfg.User == "" || cfg.Password == "" || cfg.Host == "" || cfg.Database == "" {
		return "", fmt.Errorf("invalid config provided")

	}
	return fmt.Sprintf("mysql://%s:%s@%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Database), nil
}

//	func (m *MySQLManager) WhereString(c *Criteria) (string, error) {
//		return "", nil
//	}
//
//	func (m *MySQLManager) OrderString(c *Criteria) (string, error) {
//		return "", nil
//	}
func (m *MySQLManager) LimitString(c *Criteria) string {
	if c == nil || c.Limit < 1 {
		return ""
	}
	return fmt.Sprintf(" LIMIT %d", c.Limit)
}
func (m *MySQLManager) OffsetString(c *Criteria) string {
	if c == nil || c.Offset < 1 {
		return ""
	}
	return fmt.Sprintf(" OFFSET %d", c.Offset)
}

func (m *MySQLManager) IdentityString(f string) string {
	return fmt.Sprintf("`%s`", f)
}

func (m *MySQLManager) TableCreate() string {
	return "CREATE TABLE IF NOT EXISTS `%s` (%s);"
}

func (m *MySQLManager) IndexCreate() string {
	return "CREATE INDEX `%s_%s_Idx` ON %s(`%s`);"

}
func (m *MySQLManager) TableDefinition(md Modeller) ([]string, bool) {
	return nil, false
}
func (m *MySQLManager) InsertCommand(md Modeller) string {
	return ""
}
func (m *MySQLManager) UpdateCommand(md Modeller) string {
	return ""
}
func (m *MySQLManager) DeleteCommand(md Modeller) string {
	return ""
}
func (m *MySQLManager) RefreshCommand(md Modeller) string {
	return ""
}

func (m *MySQLManager) BuildQuery(where string, order string, limit string, offset string) string {
	res := ""
	if where != "" {
		res += fmt.Sprintf(" WHERE %s", where)
	}
	if order != "" {
		res += fmt.Sprintf(" ORDER BY %s", order)
	}
	if limit != "" {
		res += fmt.Sprintf(" LIMIT %s", limit)
	}
	if offset != "" {
		res += fmt.Sprintf(" OFFSET %s", offset)
	}
	return res
}
func (m *MySQLManager) TableTest(mdl Modeller) ([]field, string, bool) {
	return nil, "", true
}

func (m *MySQLManager) TableExistsQuery(dbName, name string) string {
	return fmt.Sprintf("SHOW TABLES WHERE Tables_in_%s = '%s'", dbName, name)
}

func (m *MySQLManager) MassDelete(mdl Modeller, c *Criteria) string {
	name := getTableName(mdl)
	s := fmt.Sprintf("DELETE FROM %s", m.IdentityString(name))
	whereAdded := false
	if c != nil && c.Where != "" {
		s += fmt.Sprintf(" WHERE %s", c.WhereString(m))
		whereAdded = true
	}
	if whereAdded {
		s += " AND `DeleteDate` IS NULL"
	} else {
		s += " WHERE `DeleteDate` IS NULL"
	}
	return s
}
func (m *MySQLManager) MassDisable(mdl Modeller, c *Criteria) string {
	tm := time.Now()
	name := getTableName(mdl)
	s := fmt.Sprintf("UPDATE %s SET `DeleteDate` = '%v'", m.IdentityString(name), utils.TimeToSQL(tm))
	whereAdded := false
	if c != nil && c.Where != "" {
		s += fmt.Sprintf(" WHERE %s", c.WhereString(m))
		whereAdded = true
	}
	if whereAdded {
		s += " AND `DeleteDate` IS NULL"
	} else {
		s += " WHERE `DeleteDate` IS NULL"
	}
	return s
}

func (m *MySQLManager) Operators() []string {
	return []string{
		"`%s` = %s",
		"`%s` > %s",
		"`%s` < %s",
		"`%s` LIKE %s",
		"`%s` IN (%s)",
		"`%s` BETWEEN %s AND %s",
		"`%s` IS NULL",
		"`%s` <> %s",
		"`%s` <= %s",
		"`%s` >= %s",
		"`%s` NOT LIKE %s",
		"`%s` NOT IN (%s)",
		"`%s` NOT BETWEEN %s AND %s",
		"`%s` IS NOT NULL",
	}
}
