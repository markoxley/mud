// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
package dtorm

import (
	"database/sql"
	"errors"
	"fmt"
	"iter"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/markoxley/dtorm/utils"
	"github.com/markoxley/dtorm/where"
	uuid "github.com/satori/go.uuid"
)

type DB struct {
	cfg              *Config
	connectionString string
	dbtype           string
	mgr              Manager
	db               *sql.DB
	knownTables      []string
	tableDef         map[string][]field
}

func New(config *Config) (*DB, error) {
	mgr, err := GetManager(config)
	if err != nil {
		return nil, err
	}
	cs, err := mgr.ConnectionString(config)
	if err != nil {
		return nil, err
	}
	db := &DB{
		mgr:              mgr,
		cfg:              config,
		dbtype:           config.Type,
		connectionString: cs,
		knownTables:      make([]string, 0),
		tableDef:         make(map[string][]field),
	}
	svr, err := db.connect()
	if err != nil {
		return nil, err
	}
	db.db = svr
	mgr.SetDB(db)
	return db, nil
}

func (db *DB) Close() {
	if db.db != nil {
		db.disconnect(db.db)
	}
}

// connect attempts to connect to the database
// @return *sql.DB
// @return error
func (db *DB) connect() (*sql.DB, error) {
	tdb, err := sql.Open(db.dbtype, db.connectionString)
	if err != nil {
		return nil, err
	}
	return tdb, nil
}

// disconnect from the database
// @param db
func (db *DB) disconnect(d *sql.DB) {
	if d != nil {
		d.Close()
	}
}

// beginTransaction begins the transaction process
// @param db
// @return *sql.Tx
// @return error
func (db *DB) beginTransaction(d *sql.DB) (*sql.Tx, error) {
	return d.Begin()
}

// commitTransaction commites the transaction to the database
// @param tx
func (db *DB) CommitTransaction(tx *sql.Tx) error {
	if tx != nil {
		return tx.Commit()
	}
	return nil
}

func (db *DB) RollbackTransaction(tx *sql.Tx) error {
	if tx != nil {
		return tx.Rollback()
	}
	return nil
}

func (db *DB) BeginTransaction() (*sql.Tx, error) {
	return db.beginTransaction(db.db)
}

// selectScalar atempts to execute the specified query and returns
// the value of the first column of the first row
// @param q
// @return interface{}
// @return bool
func (db *DB) selectScalar(q string, tx ...*sql.Tx) (interface{}, bool) {
	var qtx *sql.Tx
	var err error
	if len(tx) > 0 {
		qtx = tx[0]
	} else if !db.cfg.DisabledTransactions {
		qtx, err = db.beginTransaction(db.db)
		if err != nil {
			return nil, false
		}
		defer db.CommitTransaction(qtx)
	}

	var res *sql.Rows
	if qtx != nil {
		res, err = qtx.Query(q)
	} else {
		res, err = db.db.Query(q)
	}
	if err != nil {
		return nil, false
	}
	defer res.Close()
	if res.Next() {
		var cols string
		vl := &cols
		//var vl interface{}
		res.Scan(vl)
		return cols, true
	}
	return nil, false
}

// selectQuery attempts to execute the query passed, returning
// a slice of the type specified by the type parameter
// @param q
// @return []*T
// @return bool
func (db *DB) selectQuery(m Modeller, q string, tx ...*sql.Tx) ([]Modeller, bool) {
	var qtx *sql.Tx
	var err error
	if len(tx) > 0 {
		qtx = tx[0]
	} else if !db.cfg.DisabledTransactions {
		qtx, err = db.beginTransaction(db.db)
		if err != nil {
			return nil, false
		}
		defer db.CommitTransaction(qtx)
	}
	var res *sql.Rows
	if qtx != nil {
		res, err = qtx.Query(q)
	} else {
		res, err = db.db.Query(q)
	}
	if err != nil {
		return nil, false
	}
	defer res.Close()
	return db.populateModel(m, res)
}

// populateModel creates a new slice of models of the type
// specified by the type parameter and populates the fields from the sql query
// @param r
// @return []*T
// @return bool
func (db *DB) populateModel(m Modeller, r *sql.Rows) ([]Modeller, bool) {
	s := reflect.TypeOf(m)

	res := make([]Modeller, 0, 100)
	ok := true

	// Get the column count
	cc, _ := r.Columns()

	// Make them all uppercase
	for i := range cc {
		cc[i] = strings.ToUpper(cc[i])
	}

	// Get the fields of the model and build a map of them
	//t := reflect.TypeOf(*m)
	mdl := reflect.New(s).Interface().(Modeller)
	flds, ok := db.tableDef[getTableName(mdl)]
	if !ok {
		return nil, false
	}
	fMap := make(map[string]field, len(flds))
	for _, f := range flds {
		fMap[strings.ToUpper(f.name)] = f
	}

	cols := make([]*string, len(cc))
	vls := make([]interface{}, len(cc))

	rowCount := 0
	for r.Next() {
		// s := reflect.TypeOf(m)
		v := reflect.New(s)

		for i := range cols {
			vls[i] = &cols[i]
		}
		r.Scan(vls...)

		for i := 0; i < len(cc); i++ {
			if cols[i] == nil {
				continue
			}
			if cc[i] == "ID" {
				tmpID := cols[i]
				v.Elem().FieldByName("ID").Set(reflect.ValueOf(tmpID))
			} else if cc[i] == "CREATEDATE" {
				if cols[i] != nil {
					if tm, ok := utils.SQLToTime(*cols[i]); ok {
						tmpCreate := *tm
						v.Elem().FieldByName("CreateDate").Set(reflect.ValueOf(tmpCreate))
					}
				}
			} else if cc[i] == "LASTUPDATE" {
				if cols[i] != nil {
					if tm, ok := utils.SQLToTime(*cols[i]); ok {
						tmpUpdate := *tm
						v.Elem().FieldByName("LastUpdate").Set(reflect.ValueOf(tmpUpdate))
					}
				}
			} else if cc[i] == "DELETEDATE" {
				if cols[i] != nil {
					if tm, ok := utils.SQLToTime(*cols[i]); ok {
						tmpDeleted := tm
						v.Elem().FieldByName("DeleteDate").Set(reflect.ValueOf(tmpDeleted))
					}
				}
			} else if fld, ok := fMap[cc[i]]; ok {
				switch fld.fType {
				case tInt, tLong:
					if fld.unsigned {
						if val, err := strconv.ParseUint(*cols[i], 10, 64); err != nil {
							if fld.allowNull {
								v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(val))
							} else {
								v.Elem().FieldByName(fld.name).SetUint(val)
							}
						}
					} else {
						if val, err := strconv.ParseInt(*cols[i], 10, 64); err == nil {
							if fld.allowNull {
								v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(val))
							} else {
								v.Elem().FieldByName(fld.name).SetInt(val)
							}
						}
					}
				case tBool:
					if val, err := strconv.ParseInt(*cols[i], 10, 64); err == nil {
						if fld.allowNull {
							v.Elem().FieldByName(fld.name).Elem().SetBool(val == 1)
						} else {
							v.Elem().FieldByName(fld.name).SetBool(val == 1)
						}
					}
				case tDecimal, tFloat, tDouble:
					if val, err := strconv.ParseFloat(*cols[i], 64); err == nil {
						if fld.allowNull {
							v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(val))
						} else {
							v.Elem().FieldByName(fld.name).SetFloat(val)
						}
					}
				case tDateTime:
					if cols[i] != nil {
						if val, ok := utils.SQLToTime(*cols[i]); ok {
							if fld.allowNull {
								v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(val))
							} else {
								v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(*val))
							}
						}
					}
				case tChar:
					st := *cols[i]
					strVal := st[:1]
					if fld.allowNull {
						v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(&strVal))
					} else {
						v.Elem().FieldByName(fld.name).SetString(strVal)
					}
				case tString:
					if fld.allowNull {
						strVal := *cols[i]
						v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(&strVal))
					} else {
						v.Elem().FieldByName(fld.name).SetString(*cols[i])
					}
				case tUUID:
					if fld.allowNull {
						strVal := *cols[i]
						v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(&strVal))
					} else {
						v.Elem().FieldByName(fld.name).SetString(*cols[i])
					}
				}
			}

		}
		newObj := v.Elem().Interface().(Modeller)
		db.doRestore(newObj)
		res = append(res, newObj)
		rowCount++

	}
	return res[:rowCount], ok
}

func (db *DB) populateSingle(m Modeller, r *sql.Rows) bool {
	s := reflect.TypeOf(m)
	ok := true

	// Get the column count
	cc, _ := r.Columns()

	// Make them all uppercase
	for i := range cc {
		cc[i] = strings.ToUpper(cc[i])
	}

	// Get the fields of the model and build a map of them
	//t := reflect.TypeOf(*m)
	mdl := reflect.New(s).Interface().(Modeller)
	flds, ok := db.tableDef[getTableName(mdl)]
	if !ok {
		return false
	}
	fMap := make(map[string]field, len(flds))
	for _, f := range flds {
		fMap[strings.ToUpper(f.name)] = f
	}

	cols := make([]*string, len(cc))
	vls := make([]interface{}, len(cc))

	v := reflect.New(s)

	for i := range cols {
		vls[i] = &cols[i]
	}
	r.Scan(vls...)

	for i := 0; i < len(cc); i++ {
		if cols[i] == nil {
			continue
		}
		if cc[i] == "ID" {
			tmpID := cols[i]
			v.Elem().FieldByName("ID").Set(reflect.ValueOf(tmpID))
		} else if cc[i] == "CREATEDATE" {
			if cols[i] != nil {
				if tm, ok := utils.SQLToTime(*cols[i]); ok {
					tmpCreate := *tm
					v.Elem().FieldByName("CreateDate").Set(reflect.ValueOf(tmpCreate))
				}
			}
		} else if cc[i] == "LASTUPDATE" {
			if cols[i] != nil {
				if tm, ok := utils.SQLToTime(*cols[i]); ok {
					tmpUpdate := *tm
					v.Elem().FieldByName("LastUpdate").Set(reflect.ValueOf(tmpUpdate))
				}
			}
		} else if cc[i] == "DELETEDATE" {
			if cols[i] != nil {
				if tm, ok := utils.SQLToTime(*cols[i]); ok {
					tmpDeleted := tm
					v.Elem().FieldByName("DeleteDate").Set(reflect.ValueOf(tmpDeleted))
				}
			}
		} else if fld, ok := fMap[cc[i]]; ok {
			switch fld.fType {
			case tInt, tLong:
				if fld.unsigned {
					if val, err := strconv.ParseUint(*cols[i], 10, 64); err != nil {
						if fld.allowNull {
							v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(val))
						} else {
							v.Elem().FieldByName(fld.name).SetUint(val)
						}
					}
				} else {
					if val, err := strconv.ParseInt(*cols[i], 10, 64); err == nil {
						if fld.allowNull {
							v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(val))
						} else {
							v.Elem().FieldByName(fld.name).SetInt(val)
						}
					}
				}
			case tBool:
				if val, err := strconv.ParseInt(*cols[i], 10, 64); err == nil {
					if fld.allowNull {
						v.Elem().FieldByName(fld.name).Elem().SetBool(val == 1)
					} else {
						v.Elem().FieldByName(fld.name).SetBool(val == 1)
					}
				}
			case tDecimal, tFloat, tDouble:
				if val, err := strconv.ParseFloat(*cols[i], 64); err == nil {
					if fld.allowNull {
						v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(val))
					} else {
						v.Elem().FieldByName(fld.name).SetFloat(val)
					}
				}
			case tDateTime:
				if cols[i] != nil {
					if val, ok := utils.SQLToTime(*cols[i]); ok {
						if fld.allowNull {
							v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(val))
						} else {
							v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(*val))
						}
					}
				}
			case tChar:
				st := *cols[i]
				strVal := st[:1]
				if fld.allowNull {
					v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(&strVal))
				} else {
					v.Elem().FieldByName(fld.name).SetString(strVal)
				}
			case tString:
				if fld.allowNull {
					strVal := *cols[i]
					v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(&strVal))
				} else {
					v.Elem().FieldByName(fld.name).SetString(*cols[i])
				}
			case tUUID:
				if fld.allowNull {
					strVal := *cols[i]
					v.Elem().FieldByName(fld.name).Set(reflect.ValueOf(&strVal))
				} else {
					v.Elem().FieldByName(fld.name).SetString(*cols[i])
				}
			}
		}

	}
	newObj := v.Elem().Interface().(Modeller)
	db.doRestore(newObj)
	mdl = newObj
	//	}
	return ok
}
func (db *DB) doRestore(m Modeller) {
	if r, ok := m.(Restorer); ok {
		r.Restore(db.mgr)
	}
}

// updateModel updates the date fields of the specified model
// @param m
// @param id
// @param createdate
// @param updatedate
// @param deletedate
func (db *DB) updateModel(m Modeller, id string, createdate time.Time, updatedate time.Time, deletedate *time.Time) {
	v := reflect.ValueOf(m)
	createdateValue := reflect.ValueOf(createdate)
	updatedateValue := reflect.ValueOf(updatedate)
	deletedateValue := reflect.ValueOf(deletedate)
	rv := reflect.New(reflect.TypeOf(id))
	rv.Elem().Set(reflect.ValueOf(id))

	v.Elem().FieldByName("ID").Set(rv)
	v.Elem().FieldByName("CreateDate").Set(createdateValue)
	v.Elem().FieldByName("LastUpdate").Set(updatedateValue)
	v.Elem().FieldByName("DeleteDate").Set(deletedateValue)
}

// executeQuery attempts to execute the passed sql query
// @param q
// @return bool
func (db *DB) executeQuery(q string, tx ...*sql.Tx) error {
	var qtx *sql.Tx
	var err error
	if len(tx) > 0 {
		qtx = tx[0]
	} else if !db.cfg.DisabledTransactions {
		qtx, err = db.beginTransaction(db.db)
		defer db.CommitTransaction(qtx)
	}
	if err != nil {
		return err
	}
	if qtx != nil {
		_, err = qtx.Exec(q)
	} else {
		_, err = db.db.Exec(q)
	}
	return err
}

// tableExists tests for the existence of the specified table
// @param t
// @return bool
func (db *DB) tableExists(t string) bool {
	if slices.Contains(db.knownTables, t) {
		return true
	}

	qry := db.mgr.TableExistsQuery(t)
	if _, ok := db.selectScalar(qry); ok {
		db.knownTables = append(db.knownTables, t)
		return true
	}
	return false
	// if _, ok := selectScalar(fmt.Sprintf("SHOW TABLES WHERE Tables_in_%s = '%s'", conf.Name, t)); ok {
	// 	knownTables = append(knownTables, t)
	// 	return true
	// }
	// return false
}

// RawExecute executes a sql statement on the database, without returning a value
// Not recommended for general use - can break shadowing
// @param sql
// @return bool
func (db *DB) RawExecute(sql string, tx ...*sql.Tx) error {
	return db.executeQuery(sql, tx...)
}

// RawScalar exeutes a raw sql statement that returns a single value
// Not recommended for general use
// @param sql
// @return interface{}
// @return bool
func (db *DB) RawScalar(sql string, tx ...*sql.Tx) (interface{}, bool) {
	return db.selectScalar(sql, tx...)
}

// RawSelect executes a raw sql statement on the database
// Not recommended for general use
// @param sql
// @return []map
func (db *DB) RawSelect(qry string, tx ...*sql.Tx) ([]map[string]interface{}, error) {
	var qtx *sql.Tx
	var err error
	if len(tx) > 0 {
		qtx = tx[0]
	} else if !db.cfg.DisabledTransactions {
		qtx, err = db.beginTransaction(db.db)
		if err != nil {
			return nil, err
		}
		defer db.CommitTransaction(qtx)
	}
	var res *sql.Rows
	if qtx != nil {
		res, err = qtx.Query(qry)
	} else {
		res, err = db.db.Query(qry)
	}
	if err != nil {
		return nil, err
	}
	defer res.Close()
	data := make([]map[string]interface{}, 0, 10)

	// Get the column count
	cc, _ := res.Columns()

	cols := make([]*string, len(cc))
	vls := make([]interface{}, len(cc))

	for res.Next() {

		for i := range cols {
			vls[i] = &cols[i]
		}
		res.Scan(vls...)
		row := make(map[string]interface{})
		for i, n := range cc {
			tmp := vls[i]
			row[n] = *(tmp.(**string)) //vls[i]
		}
		data = append(data, row)
	}
	return data, nil
}

// getCriteria returns the criteria for a query in SQL format
// @param criteria
// @return *Criteria
// @return error
func (db *DB) getCriteria(criteria []interface{}) (*Criteria, error) {
	for _, cr := range criteria {
		if cr == nil {
			continue
		}

		if c, ok := cr.(*Criteria); ok {
			return c, nil
		} else if c, ok := cr.(Criteria); ok {
			return &c, nil
		} else if c, ok := cr.(string); ok {
			var re = regexp.MustCompile(`^\s*[0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12}\s*$`)
			if len(re.FindStringIndex(c)) > 0 {
				return &Criteria{Where: where.Equal("id", c)}, nil
			}

			return &Criteria{Where: c}, nil
		} else if c, ok := cr.(fmt.Stringer); ok {
			return &Criteria{Where: c.String()}, nil
		}

		return nil, errors.New("invalid criteria format")
	}
	return &Criteria{}, nil
}

func (db *DB) Range(mdl Modeller, criteria ...interface{}) iter.Seq[Modeller] {
	return func(yield func(Modeller) bool) {
		c, err := db.getCriteria(criteria)
		if err != nil {
			return
		}
		t := reflect.TypeOf(mdl)
		n := t.Name()
		if !db.tableExists(getTableName(mdl)) {
			return
		}
		s := fmt.Sprintf("select * from %s", db.mgr.IdentityString(n))
		s += c.String(db.mgr)

		var qtx *sql.Tx
		if !db.cfg.DisabledTransactions {
			qtx, err = db.beginTransaction(db.db)
			if err != nil {
				return
			}
			defer db.CommitTransaction(qtx)
		}
		var res *sql.Rows
		if qtx != nil {
			res, err = qtx.Query(s)
		} else {
			res, err = db.db.Query(s)
		}
		if err != nil {
			return
		}
		defer res.Close()
		for res.Next() {
			r := reflect.New(reflect.TypeOf(mdl).Elem()).Interface().(Modeller)
			db.populateSingle(r, res)
			if !yield(r) {
				return
			}
		}
	}
}

// Fetch populates the slice with models from the database that match the criteria.
// Returns an error if this fails
// @param criteria
// @return []*T
// @return error
func (db *DB) Fetch(mdl Modeller, criteria ...interface{}) ([]Modeller, error) {
	c, err := db.getCriteria(criteria)
	if err != nil {
		return nil, err
	}

	t := reflect.TypeOf(mdl)
	n := t.Name()
	if !db.tableExists(getTableName(mdl)) {
		return nil, errors.New("failed table check")
	}

	s := fmt.Sprintf("SELECT * FROM %s", db.mgr.IdentityString(n))
	s += c.String(db.mgr)
	res, ok := db.selectQuery(mdl, s)
	if !ok {
		return nil, errors.New("error selecting data")
	}
	return res, nil
}

// First returns the first model that matches the criteria
// @param criteria
// @return *T
// @return error
func (db *DB) First(m Modeller, criteria ...interface{}) error {
	c, err := db.getCriteria(criteria)
	if err != nil {
		return err
	}
	c.Limit = 1
	c.Offset = 0
	ml, err := db.Fetch(m, c)
	if err != nil {
		return err
	}
	if len(ml) == 0 {
		return NoResults("no results")
	}
	reflect.ValueOf(m).Elem().Set(reflect.ValueOf(ml[0]).Elem())
	return nil
}

// Count returns the number of rows in the database that match the criteria
// @param criteria
// @return int
func (db *DB) Count(m Modeller, criteria ...interface{}) int {
	c, err := db.getCriteria(criteria)
	if err != nil {
		return -1
	}
	t := getTableName(m)
	if !db.tableExists(t) {
		return 0
	}
	s := fmt.Sprintf("Select Count(*) from %s", db.mgr.IdentityString(t))
	if c != nil {
		s += c.WhereString(db.mgr)
	}
	if i, ok := db.selectScalar(s); ok {
		if vl, vlok := i.(string); vlok {
			if res, err := strconv.Atoi(vl); err == nil {
				return res
			}
		}
	}
	return 0

}

func (db *DB) insertCommand(m Modeller) (string, error) {
	t := reflect.TypeOf(m)
	n := t.Name()
	if !db.tableExists(getTableName(m)) {
		return "", errors.New("failed table check")
	}
	flds, n, err := db.tableTest(m)
	if err != nil {
		return "", err
	}
	uid := uuid.NewV4()

	fds := "ID, CreateDate, LastUpdate"
	now := time.Now()
	db.updateModel(m, uid.String(), now, now, nil)
	q := fmt.Sprintf("'%s', '%s', '%s'", uid, now, now)
	v := reflect.ValueOf(m).Elem()
	for _, f := range flds {
		if f.name == "ID" || f.name == "CreateDate" || f.name == "LastUpdate" || f.name == "DeleteDate" {
			continue
		}
		vi := v.FieldByName(f.name)

		if f.allowNull {
			if vi.IsNil() {
				continue
			}
			vi = vi.Elem()
		}

		vf := vi.Interface()

		if vl, ok := utils.MakeValue(vf); ok {
			fds += fmt.Sprintf(", %s", db.mgr.IdentityString(f.name))
			q += fmt.Sprintf(", %s", vl)
		}
	}

	def := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", n, fds, q)
	return def, nil
}

// updateCommand returns the SQL command to update the
// current model in the database
// @param m
// @return string
func (db *DB) updateCommand(m Modeller) (string, error) {
	flds, n, err := db.tableTest(m)
	if err != nil {
		return "", err
	}
	now := time.Now()
	db.updateLastUpdate(m, now)
	res := fmt.Sprintf("UPDATE %s SET", n)
	v := reflect.ValueOf(m)
	first := true
	for _, f := range flds {
		if f.name != "ID" && f.name != "CreateDate" {
			if !first {
				res += ","
			}
			first = false
			var value interface{}
			if f.allowNull {
				if v.Elem().FieldByName(f.name).IsNil() {
					res += fmt.Sprintf(" %s = null", db.mgr.IdentityString(f.name))
					continue
				}
				value = v.Elem().FieldByName(f.name).Elem().Interface()
			} else {
				value = v.Elem().FieldByName(f.name).Interface()
			}
			if vl, ok := utils.MakeValue(value); ok {
				res += fmt.Sprintf(" %s = %s", db.mgr.IdentityString(f.name), vl)
			}
		}
	}
	def := res + fmt.Sprintf(" WHERE %s = '%s'", db.mgr.IdentityString("Id"), *m.GetID())
	return def, nil
}

func (db *DB) updateLastUpdate(m Modeller, date time.Time) {
	v := reflect.ValueOf(m)
	dateValue := reflect.ValueOf(date)
	v.Elem().FieldByName("LastUpdate").Set(dateValue)
}
func (db *DB) tableTest(m Modeller) ([]field, string, error) {
	n := getTableName(m)
	sql, reqd := db.tableDefinition(m)
	if reqd {
		te := db.tableExists(n)
		db.knownTables = append(db.knownTables, n)
		if !te {
			for _, s := range sql {
				if err := db.executeQuery(s); err != nil {
					return nil, "", err
				}
			}
			if sd := m.StandingData(); sd != nil {
				for _, data := range sd {
					db.Save(data)
				}
			}
		}
	}
	flds, ok := db.tableDef[n]
	if !ok {
		return nil, "", errors.New("table definition not found")
	}
	return flds, n, nil
}

// Save stores the model in the database.
// Depending on the status of the model, this is either
// an update or an insert command
// @param m
// @return bool
func (db *DB) Save(m Modeller) bool {
	if u, ok := m.(Updater); ok {
		cmd, err := u.Update(db.mgr)
		if err != nil {
			return false
		}
		return db.executeQuery(cmd) != nil
	}
	if m.IsNew() {
		cmd, err := db.insertCommand(m)
		if err != nil {
			return false
		}
		return db.executeQuery(cmd) != nil
	}
	updCmd, err := db.updateCommand(m)
	if err != nil {
		return false
	}
	return db.executeQuery(updCmd) != nil
}

// Remove removes the passed model from the database
// @param m
// @return bool
func (db *DB) Remove(m Modeller) bool {
	if m.GetID() == nil {
		return true
	}
	c := &Criteria{
		Where: where.Equal("ID", m.GetID()),
	}
	var s string
	if db.cfg.Deletable {
		s = db.massDelete(m, c)
	} else {
		s = db.massDisable(m, c)
	}
	return db.executeQuery(s) != nil
}

func (db *DB) massDelete(m Modeller, c *Criteria) string {
	name := getTableName(m)
	s := fmt.Sprintf("DELETE FROM %s", db.mgr.IdentityString(name))
	whereAdded := false
	if c != nil && c.Where != "" {
		s += fmt.Sprintf(" WHERE %s", c.WhereString(db.mgr))
		whereAdded = true
	}
	if whereAdded {
		s += fmt.Sprintf(" AND %s IS NULL", db.mgr.IdentityString("[DeleteDate]"))
	} else {
		s += fmt.Sprintf(" WHERE %s IS NULL", db.mgr.IdentityString("[DeleteDate]"))
	}
	return s
}

func (db *DB) massDisable(m Modeller, c *Criteria) string {
	tm := time.Now()
	deleteDate := db.mgr.IdentityString("[DeleteDate]")
	name := getTableName(m)
	s := fmt.Sprintf("UPDATE %s SET %s = '%v'", db.mgr.IdentityString(name), deleteDate, utils.TimeToSQL(tm))
	whereAdded := false
	if c != nil && c.Where != "" {
		s += fmt.Sprintf(" WHERE %s", c.WhereString(db.mgr))
		whereAdded = true
	}
	if whereAdded {
		s += fmt.Sprintf(" AND %s IS NULL", deleteDate)
	} else {
		s += fmt.Sprintf(" WHERE %s IS NULL", deleteDate)
	}
	return s
}

// RemoveMany removes all models of the specified type that match the criteria
// @param c
// @return int
// @return bool
func (db *DB) RemoveMany(m Modeller, c *Criteria) (int, bool) {
	t := getTableName(m)
	if !db.tableExists(t) {
		return 0, true
	}
	r := db.Count(m, c)
	if r == 0 {
		return 0, true
	}
	s := ""
	if db.cfg.Deletable {
		s = db.massDelete(m, c)
	} else {
		s = db.massDisable(m, c)
	}
	ok := db.executeQuery(s) != nil
	return r, ok
}

func (db *DB) tableDefinition(m Modeller) ([]string, bool) {
	sql := make([]string, 0, 3)

	n := getTableName(m)
	if _, ok := db.tableDef[n]; ok {
		return nil, false
	}

	t := reflect.TypeOf(m)
	var nm interface{}
	if t.Kind() == reflect.Ptr {
		nm = reflect.New(t.Elem()).Elem().Interface()
	} else {
		nm = reflect.New(t).Elem().Interface()
	}
	fs := getDefs(nm, true)

	db.tableDef[n] = fs
	if len(fs) == 0 {
		return nil, false
	}
	flds := ""
	keys := make([]string, 0, 5)
	for _, f := range fs {
		if flds != "" {
			flds += ", "
		}
		flds += fmt.Sprintf("`%s` %s", f.name, pgFieldNames[f.fType])
		if f.fType != tUUID && f.fType != tChar && f.size.size > 0 {
			flds += fmt.Sprintf("(%s)", f.size.String())
		}
		if f.fType == tString && f.size.size == 0 {
			flds += "(256)"
		}
		if f.unsigned {
			flds += " UNSIGNED"
		}
		if !f.allowNull {
			flds += " NOT NULL"
		}
		if f.key {
			keys = append(keys, f.name)
		}
	}
	sql = append(sql, fmt.Sprintf(db.mgr.TableCreate(), n, flds))
	kn := strings.ReplaceAll(n, ".", "_")
	for _, k := range keys {
		sql = append(sql, fmt.Sprintf(db.mgr.IndexCreate(), kn, k, n, k))
	}
	return sql, true
}

func (db *DB) Refresh(m Modeller) error {
	if m.GetID() == nil {
		return errors.New("no id")
	}
	return db.First(m, where.Equal("ID", m.GetID()))
}

