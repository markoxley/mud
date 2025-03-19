package dtorm

import (
	"reflect"
	"time"

	"github.com/markoxley/dtorm/utils"
)

type Model struct {
	ID         *string
	CreateDate time.Time
	LastUpdate time.Time
	DeleteDate *time.Time
	tableName  *string
}

func CreateModel() Model {
	return Model{
		CreateDate: time.Now(),
		LastUpdate: time.Now(),
	}
}

func (m Model) StandingData() []Modeller {
	return nil
}

func (m Model) GetID() *string {
	return m.ID
}

func (m Model) IsNew() bool {
	return m.ID == nil
}

func (m Model) IsDeleted() bool {
	return m.DeleteDate == nil
}

func (m *Model) Disable() {
	m.DeleteDate = utils.Ptr(time.Now())
}
func getTableName(m Modeller) string {
	if reflect.TypeOf(m).Kind() == reflect.Pointer {
		return reflect.Indirect(reflect.ValueOf(m).Elem()).Type().Name()
	}
	return reflect.ValueOf(m).Type().Name()

}
