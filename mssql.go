package dtorm

import "fmt"

type MSSQLManager struct {
	// Needs to implement Manager
}

func (m *MSSQLManager) ConnectionString(cfg *Config) (string, error) {
	if cfg == nil {
		return "", fmt.Errorf("no config provided")

	}
	if cfg.User == "" || cfg.Password == "" || cfg.Host == "" || cfg.Database == "" {
		return "", fmt.Errorf("invalid config provided")

	}
	return fmt.Sprintf("sqlserver://%s:%s@%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Database), nil
}

func (m *MSSQLManager) WhereString(c *Criteria) (string, error) {
	return "", nil
}
func (m *MSSQLManager) OrderString(c *Criteria) (string, error) {
	return "", nil
}
func (m *MSSQLManager) LimitString(c *Criteria) (string, error) {
	return "", nil
}
func (m *MSSQLManager) OffsetString(c *Criteria) (string, error) {
	return "", nil
}
func (m *MSSQLManager) IdentityString(f string) string {
	return fmt.Sprintf("[%s]", f)
}

func (m *MSSQLManager) TableDefinition(md Modeller) ([]string, bool) {
	return nil, false
}
func (m *MSSQLManager) InsertCommand(md Modeller) string {
	return ""
}
func (m *MSSQLManager) UpdateCommand(md Modeller) string {
	return ""
}
func (m *MSSQLManager) DeleteCommand(md Modeller) string {
	return ""
}
func (m *MSSQLManager) RefreshCommand(md Modeller) string {
	return ""
}
