package dtorm

import "fmt"

type SqliteManager struct {
	// Needs to implement Manager
}

func (m *SqliteManager) ConnectionString(cfg *Config) (string, error) {
	if cfg == nil {
		return "", fmt.Errorf("no config provided")

	}
	if cfg.Database == "" {
		return "", fmt.Errorf("invalid config provided")

	}
	return fmt.Sprintf("sqlite3://%s", cfg.Database), nil
}

func (m *SqliteManager) WhereString(c *Criteria) (string, error) {
	return "", nil
}
func (m *SqliteManager) OrderString(c *Criteria) (string, error) {
	return "", nil
}
func (m *SqliteManager) LimitString(c *Criteria) (string, error) {
	return "", nil
}
func (m *SqliteManager) OffsetString(c *Criteria) (string, error) {
	return "", nil
}
func (m *SqliteManager) IdentityString(f string) string {
	return fmt.Sprintf("\"%s\"", f)
}

func (m *SqliteManager) TableDefinition(md Modeller) ([]string, bool) {
	return nil, false
}
func (m *SqliteManager) InsertCommand(md Modeller) string {
	return ""
}
func (m *SqliteManager) UpdateCommand(md Modeller) string {
	return ""
}
func (m *SqliteManager) DeleteCommand(md Modeller) string {
	return ""
}
func (m *SqliteManager) RefreshCommand(md Modeller) string {
	return ""
}
