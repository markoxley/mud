// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
package mud

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
	"sync"
	"time"

	"github.com/markoxley/mud/order"
	"github.com/markoxley/mud/utils"
	"github.com/markoxley/mud/where"
	uuid "github.com/satori/go.uuid"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/microsoft/go-mssqldb"
	_ "modernc.org/sqlite"
)

type DB struct {
	cfg              *Config
	connectionString string
	dbtype           string
	mgr              Manager
	db               *sql.DB
	mu               sync.RWMutex
	// createMu serializes table creation: the definition is cached before
	// the CREATE statement runs, so concurrent callers must wait for it.
	createMu    sync.Mutex
	knownTables []string
	tableDef    map[string][]field
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
	var dbType string
	switch config.Type {
	case "sqlite", "sqlite3":
		dbType = "sqlite"
	case "mysql", "mariadb":
		dbType = "mysql"
	case "sqlserver", "mssql":
		dbType = "sqlserver"
	default:
		return nil, fmt.Errorf("invalid database type: %s", config.Type)
	}
	db := &DB{
		mgr:              mgr,
		cfg:              config,
		dbtype:           dbType,
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

func (db *DB) SqlDB() *sql.DB {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.db
}

func (db *DB) Close() {
	db.mu.Lock()
	defer db.mu.Unlock()
	if db.db != nil {
		db.disconnect(db.db)
		db.db = nil
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
	if db.dbtype == "sqlite" {
		return nil, nil
	}
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
// @param args
// @return interface{}
// @return bool
func (db *DB) selectScalar(q string, args []interface{}, tx ...*sql.Tx) (interface{}, bool) {
	var qtx *sql.Tx
	var err error
	if db.dbtype != "sqlite" {
		if len(tx) > 0 {
			qtx = tx[0]
		} else if !db.cfg.DisabledTransactions {
			qtx, err = db.beginTransaction(db.db)
			if err != nil {
				return nil, false
			}
			defer db.CommitTransaction(qtx)
		}
	}
	var res *sql.Rows
	if qtx != nil {
		res, err = qtx.Query(q, args...)
	} else {
		res, err = db.db.Query(q, args...)
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
// @param args
// @return []*T
// @return bool
func (db *DB) selectQuery(m Modeller, q string, args []interface{}, tx ...*sql.Tx) ([]Modeller, bool) {
	var qtx *sql.Tx
	var err error
	if db.dbtype != "sqlite" {
		if len(tx) > 0 {
			qtx = tx[0]
		} else if !db.cfg.DisabledTransactions {
			qtx, err = db.beginTransaction(db.db)
			if err != nil {
				return nil, false
			}
			defer db.CommitTransaction(qtx)
		}
	}
	var res *sql.Rows
	if qtx != nil {
		res, err = qtx.Query(q, args...)
	} else {
		res, err = db.db.Query(q, args...)
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

	// Get the column names
	cc, err := r.Columns()
	if err != nil {
		return nil, false
	}

	// Make them all uppercase
	for i := range cc {
		cc[i] = strings.ToUpper(cc[i])
	}

	// Get the fields of the model and build a map of them
	mdl := reflect.New(s).Interface().(Modeller)
	db.mu.RLock()
	flds, ok := db.tableDef[GetTableName(mdl)]
	db.mu.RUnlock()
	if !ok {
		return nil, false
	}
	fMap := fieldMap(flds)

	cols := make([]*string, len(cc))
	vls := make([]interface{}, len(cc))
	for i := range cols {
		vls[i] = &cols[i]
	}

	for r.Next() {
		v := reflect.New(s)

		if err := r.Scan(vls...); err != nil {
			return res, false
		}
		populateFields(v, cc, cols, fMap)
		newObj := v.Elem().Interface().(Modeller)
		db.doRestore(newObj)
		res = append(res, newObj)
	}
	return res, true
}

func (db *DB) populateSingle(m Modeller, r *sql.Rows) (Modeller, bool) {
	s := reflect.TypeOf(m)

	// Get the column names
	cc, err := r.Columns()
	if err != nil {
		return nil, false
	}

	// Make them all uppercase
	for i := range cc {
		cc[i] = strings.ToUpper(cc[i])
	}

	// Get the fields of the model and build a map of them
	mdl := reflect.New(s).Interface().(Modeller)
	db.mu.RLock()
	flds, ok := db.tableDef[GetTableName(mdl)]
	db.mu.RUnlock()
	if !ok {
		return nil, false
	}
	fMap := fieldMap(flds)

	cols := make([]*string, len(cc))
	vls := make([]interface{}, len(cc))
	for i := range cols {
		vls[i] = &cols[i]
	}
	if err := r.Scan(vls...); err != nil {
		return nil, false
	}

	v := reflect.New(s)
	populateFields(v, cc, cols, fMap)
	newObj := v.Elem().Interface().(Modeller)
	db.doRestore(newObj)
	return newObj, true
}

// fieldMap indexes field definitions by upper-cased column name.
func fieldMap(flds []field) map[string]field {
	fMap := make(map[string]field, len(flds))
	for _, f := range flds {
		fMap[strings.ToUpper(f.name)] = f
	}
	return fMap
}

// populateFields copies the values scanned into cols into the freshly
// allocated model pointed to by v.
func populateFields(v reflect.Value, cc []string, cols []*string, fMap map[string]field) {
	for i := 0; i < len(cc); i++ {
		if cols[i] == nil {
			continue
		}
		switch cc[i] {
		case "ID":
			v.Elem().FieldByName("ID").Set(reflect.ValueOf(cols[i]))
		case "CREATEDATE":
			if tm, ok := utils.SQLToTime(*cols[i]); ok {
				v.Elem().FieldByName("CreateDate").Set(reflect.ValueOf(*tm))
			}
		case "LASTUPDATE":
			if tm, ok := utils.SQLToTime(*cols[i]); ok {
				v.Elem().FieldByName("LastUpdate").Set(reflect.ValueOf(*tm))
			}
		case "DELETEDATE":
			if tm, ok := utils.SQLToTime(*cols[i]); ok {
				v.Elem().FieldByName("DeleteDate").Set(reflect.ValueOf(tm))
			}
		default:
			if fld, ok := fMap[cc[i]]; ok {
				setFieldValue(v.Elem(), fld, *cols[i])
			}
		}
	}
}

// setNullable assigns to the element of a pointer field, allocating the
// pointer first if it is nil.
func setNullable(fv reflect.Value, set func(reflect.Value)) {
	if fv.IsNil() {
		fv.Set(reflect.New(fv.Type().Elem()))
	}
	set(fv.Elem())
}

// setFieldValue parses the scanned string s and stores it in the field
// named by fld, allocating pointer fields when fld.allowNull is set.
func setFieldValue(e reflect.Value, fld field, s string) {
	fv := e.FieldByName(fld.name)
	if !fv.IsValid() || !fv.CanSet() {
		return
	}
	switch fld.fType {
	case tInt, tLong:
		if fld.unsigned {
			if val, err := strconv.ParseUint(s, 10, 64); err == nil {
				if fld.allowNull {
					setNullable(fv, func(e reflect.Value) { e.SetUint(val) })
				} else {
					fv.SetUint(val)
				}
			}
		} else {
			if val, err := strconv.ParseInt(s, 10, 64); err == nil {
				if fld.allowNull {
					setNullable(fv, func(e reflect.Value) { e.SetInt(val) })
				} else {
					fv.SetInt(val)
				}
			}
		}
	case tBool:
		var b, parsed bool
		if val, err := strconv.ParseInt(s, 10, 64); err == nil {
			b, parsed = val == 1, true
		} else if val, err := strconv.ParseBool(s); err == nil {
			b, parsed = val, true
		}
		if parsed {
			if fld.allowNull {
				setNullable(fv, func(e reflect.Value) { e.SetBool(b) })
			} else {
				fv.SetBool(b)
			}
		}
	case tDecimal, tFloat, tDouble:
		if val, err := strconv.ParseFloat(s, 64); err == nil {
			if fld.allowNull {
				setNullable(fv, func(e reflect.Value) { e.SetFloat(val) })
			} else {
				fv.SetFloat(val)
			}
		}
	case tDateTime:
		if val, ok := utils.SQLToTime(s); ok {
			if fld.allowNull {
				setNullable(fv, func(e reflect.Value) { e.Set(reflect.ValueOf(*val)) })
			} else {
				fv.Set(reflect.ValueOf(*val))
			}
		}
	case tChar:
		strVal := ""
		if r := []rune(s); len(r) > 0 {
			strVal = string(r[:1])
		}
		if fld.allowNull {
			setNullable(fv, func(e reflect.Value) { e.SetString(strVal) })
		} else {
			fv.SetString(strVal)
		}
	case tString, tUUID:
		if fld.allowNull {
			setNullable(fv, func(e reflect.Value) { e.SetString(s) })
		} else {
			fv.SetString(s)
		}
	}
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

	v.Elem().FieldByName("ID").Set(reflect.ValueOf(&id))
	v.Elem().FieldByName("CreateDate").Set(createdateValue)
	v.Elem().FieldByName("LastUpdate").Set(updatedateValue)
	v.Elem().FieldByName("DeleteDate").Set(deletedateValue)
}

// executeQuery attempts to execute the passed sql query
// @param q
// @param args
// @return bool
func (db *DB) executeQuery(q string, args []interface{}, tx ...*sql.Tx) error {
	var qtx *sql.Tx
	var err error
	if db.dbtype != "sqlite" {
		if len(tx) > 0 {
			qtx = tx[0]
		} else if !db.cfg.DisabledTransactions {
			qtx, err = db.beginTransaction(db.db)
			defer db.CommitTransaction(qtx)
		}
		if err != nil {
			return err
		}
	}
	if qtx != nil {
		_, err = qtx.Exec(q, args...)
	} else {
		_, err = db.db.Exec(q, args...)
	}
	return err
}

// tableExists tests for the existence of the specified table
// @param t
// @return bool
func (db *DB) tableExists(t string) bool {
	db.mu.RLock()
	if slices.Contains(db.knownTables, t) {
		db.mu.RUnlock()
		return true
	}
	db.mu.RUnlock()

	qry, qargs := db.mgr.TableExistsQuery(t)
	if _, ok := db.selectScalar(qry, qargs); ok {
		db.mu.Lock()
		if !slices.Contains(db.knownTables, t) {
			db.knownTables = append(db.knownTables, t)
		}
		db.mu.Unlock()
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
// @param args
// @return bool
func (db *DB) RawExecute(sql string, args []interface{}, tx ...*sql.Tx) error {
	return db.executeQuery(sql, args, tx...)
}

// RawScalar exeutes a raw sql statement that returns a single value
// Not recommended for general use
// @param sql
// @param args
// @return interface{}
// @return bool
func (db *DB) RawScalar(sql string, args []interface{}, tx ...*sql.Tx) (interface{}, bool) {
	return db.selectScalar(sql, args, tx...)
}

// RawSelect executes a raw sql statement on the database
// Not recommended for general use
// @param sql
// @param args
// @return []map
func (db *DB) RawSelect(qry string, args []interface{}, tx ...*sql.Tx) ([]map[string]interface{}, error) {
	var qtx *sql.Tx
	var err error
	if db.dbtype != "sqlite" {
		if len(tx) > 0 {
			qtx = tx[0]
		} else if !db.cfg.DisabledTransactions {
			qtx, err = db.beginTransaction(db.db)
			if err != nil {
				return nil, err
			}
			defer db.CommitTransaction(qtx)
		}
	}
	var res *sql.Rows
	if qtx != nil {
		res, err = qtx.Query(qry, args...)
	} else {
		res, err = db.db.Query(qry, args...)
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
		} else if c, ok := cr.(*where.Builder); ok {
			return &Criteria{Where: c}, nil
		} else if c, ok := cr.(where.Builder); ok {
			return &Criteria{Where: &c}, nil
		} else if c, ok := cr.(*order.Builder); ok {
			return &Criteria{Order: c}, nil
		} else if c, ok := cr.(order.Builder); ok {
			return &Criteria{Order: &c}, nil
		} else if c, ok := cr.(string); ok {
			var re = regexp.MustCompile(`(?i)^\s*[0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12}\s*$`)
			if len(re.FindStringIndex(c)) > 0 {
				return &Criteria{Where: where.Equal("ID", c)}, nil
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
		n := GetTableName(mdl)
		_, _, err = db.tableTest(mdl)
		if err != nil {
			return
		}
		s := fmt.Sprintf("select * from %s", db.mgr.IdentityString(n))
		crit, args := c.Build(db.mgr)
		s += crit

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
			res, err = qtx.Query(s, args...)
		} else {
			res, err = db.db.Query(s, args...)
		}
		if err != nil {
			return
		}
		defer res.Close()
		for res.Next() {
			r, ok := db.populateSingle(mdl, res)
			if !ok {
				continue
			}
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

	n := GetTableName(mdl)
	_, _, err = db.tableTest(mdl)
	if err != nil {
		return nil, err
	}

	s := fmt.Sprintf("SELECT * FROM %s", db.mgr.IdentityString(n))
	crit, args := c.Build(db.mgr)
	s += crit
	res, ok := db.selectQuery(mdl, s, args)
	if !ok {
		return nil, errors.New("error selecting data")
	}
	return res, nil
}

// First returns the first model that matches the criteria
// @param criteria
// @return *T
// @return error
func (db *DB) First(m Modeller, criteria ...interface{}) (Modeller, error) {
	c, err := db.getCriteria(criteria)
	if err != nil {
		return nil, err
	}
	c.Limit = 1
	c.Offset = 0
	ml, err := db.Fetch(m, c)
	if err != nil {
		return nil, err
	}
	if len(ml) == 0 {
		return nil, NoResults("no results")
	}
	return ml[0], nil
}

// Count returns the number of rows in the database that match the criteria
// @param criteria
// @return int
func (db *DB) Count(m Modeller, criteria ...interface{}) int {
	c, err := db.getCriteria(criteria)
	if err != nil {
		return -1
	}
	t := GetTableName(m)
	_, _, err = db.tableTest(m)
	if err != nil {
		return -1
	}
	s := fmt.Sprintf("Select Count(*) from %s", db.mgr.IdentityString(t))
	var args []interface{}
	if c != nil {
		var where string
		where, args = c.WhereClause(db.mgr)
		s += where
	}
	if i, ok := db.selectScalar(s, args); ok {
		if vl, vlok := i.(string); vlok {
			if res, err := strconv.Atoi(vl); err == nil {
				return res
			}
		}
	}
	return 0

}

func (db *DB) insertCommand(m Modeller) (string, []interface{}, error) {
	//t := reflect.TypeOf(m)
	//n := getTableName(m)
	flds, n, err := db.tableTest(m)
	if err != nil {
		return "", nil, err
	}
	uid := uuid.NewV4()

	fds := "ID, CreateDate, LastUpdate"
	now := time.Now()
	db.updateModel(m, uid.String(), now, now, nil)
	pos := 1
	q := db.mgr.Placeholder(&pos) + ", " + db.mgr.Placeholder(&pos) + ", " + db.mgr.Placeholder(&pos)
	args := []interface{}{uid.String(), utils.TimeToSQL(now), utils.TimeToSQL(now)}
	v := reflect.ValueOf(m)
	for _, f := range flds {
		if f.identity || f.name == "ID" || f.name == "CreateDate" || f.name == "LastUpdate" || f.name == "DeleteDate" {
			continue
		}
		vi := v.Elem().FieldByName(f.name)
		// Ambiguous promoted fields (same name in two embedded structs)
		// resolve to an invalid reflect.Value; skip them rather than panic.
		if !vi.IsValid() || !vi.CanInterface() {
			continue
		}

		if f.allowNull {
			if vi.IsNil() {
				continue
			}
			vi = vi.Elem()
		}

		vf := vi.Interface()

		if arg, ok := utils.MakeArg(vf); ok {
			fds += fmt.Sprintf(", %s", db.mgr.IdentityString(f.name))
			q += fmt.Sprintf(", %s", db.mgr.Placeholder(&pos))
			args = append(args, arg)
		}
	}

	def := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", db.mgr.IdentityString(n), fds, q)
	return def, args, nil
}

// updateCommand returns the parameterized SQL command to update the
// current model in the database, along with the arguments to bind
// @param m
// @return string
// @return []interface{}
func (db *DB) updateCommand(m Modeller) (string, []interface{}, error) {
	flds, n, err := db.tableTest(m)
	if err != nil {
		return "", nil, err
	}
	now := time.Now()
	db.updateLastUpdate(m, now)
	res := fmt.Sprintf("UPDATE %s SET", db.mgr.IdentityString(n))
	v := reflect.ValueOf(m)
	pos := 1
	var args []interface{}
	first := true
	for _, f := range flds {
		if f.name == "ID" || f.name == "CreateDate" || f.identity {
			continue
		}
		var assignment string
		vi := v.Elem().FieldByName(f.name)
		// Ambiguous promoted fields resolve to an invalid reflect.Value;
		// skip them rather than panic.
		if !vi.IsValid() || !vi.CanInterface() {
			continue
		}
		if f.allowNull && vi.IsNil() {
			assignment = fmt.Sprintf(" %s = null", db.mgr.IdentityString(f.name))
		} else {
			var value interface{}
			if f.allowNull {
				value = vi.Elem().Interface()
			} else {
				value = vi.Interface()
			}
			if arg, ok := utils.MakeArg(value); ok {
				assignment = fmt.Sprintf(" %s = %s", db.mgr.IdentityString(f.name), db.mgr.Placeholder(&pos))
				args = append(args, arg)
			}
		}
		if assignment == "" {
			continue
		}
		if !first {
			res += ","
		}
		first = false
		res += assignment
	}
	if first {
		return "", nil, errors.New("no updatable fields")
	}
	def := res + fmt.Sprintf(" WHERE %s = %s", db.mgr.IdentityString("ID"), db.mgr.Placeholder(&pos))
	args = append(args, *m.GetID())
	return def, args, nil
}

func (db *DB) updateLastUpdate(m Modeller, date time.Time) {
	v := reflect.ValueOf(m)
	dateValue := reflect.ValueOf(date)
	v.Elem().FieldByName("LastUpdate").Set(dateValue)
}
func (db *DB) tableTest(m Modeller) ([]field, string, error) {
	n := GetTableName(m)
	// Serialize table creation: tableDefinition caches the definition
	// before the CREATE statement runs, so without this lock a concurrent
	// caller could proceed while the table does not yet exist.
	db.createMu.Lock()
	sql, reqd := db.tableDefinition(m)
	created := false
	if reqd && !db.tableExists(n) {
		for _, s := range sql {
			if err := db.executeQuery(s, nil); err != nil {
				db.createMu.Unlock()
				return nil, "", err
			}
		}
		created = true
	}
	db.createMu.Unlock()

	// Seed standing data after releasing the creation lock: Save re-enters
	// tableTest for the standing-data models.
	if created {
		for _, data := range m.StandingData() {
			if err := db.Save(data); err != nil {
				return nil, "", err
			}
		}
	}
	db.mu.RLock()
	flds, ok := db.tableDef[n]
	db.mu.RUnlock()
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
func (db *DB) Save(m Modeller, tx ...*sql.Tx) error {
	if u, ok := m.(Updater); ok {
		err := u.Update(db.mgr)
		if err != nil {
			return err
		}
	}
	if m.IsNew() {
		cmd, args, err := db.insertCommand(m)
		if err != nil {
			return err
		}
		return db.executeQuery(cmd, args)
	}
	updCmd, args, err := db.updateCommand(m)
	if err != nil {
		return err
	}
	return db.executeQuery(updCmd, args, tx...)
}

// Remove removes the passed model from the database
// @param m
// @return bool
func (db *DB) Remove(m Modeller, tx ...*sql.Tx) error {
	if m.GetID() == nil {
		return nil
	}
	c := &Criteria{
		Where: where.Equal("ID", *m.GetID()),
	}
	var s string
	var args []interface{}
	if db.cfg.Deletable {
		s, args = db.massDelete(m, c)
	} else {
		s, args = db.massDisable(m, c)
	}
	return db.executeQuery(s, args, tx...)
}

func (db *DB) massDelete(m Modeller, c *Criteria) (string, []interface{}) {
	name := GetTableName(m)
	s := fmt.Sprintf("DELETE FROM %s", db.mgr.IdentityString(name))
	var args []interface{}
	if c != nil {
		var where string
		where, args = c.WhereClause(db.mgr)
		if trimmed := strings.TrimSpace(where); trimmed != "" {
			s += fmt.Sprintf(" %s", trimmed)
		}
	}
	return s, args
}

func (db *DB) massDisable(m Modeller, c *Criteria) (string, []interface{}) {
	tm := time.Now()
	deleteDate := db.mgr.IdentityString("DeleteDate")
	name := GetTableName(m)
	pos := 2 // position 1 is used by the SET placeholder below
	s := fmt.Sprintf("UPDATE %s SET %s = %s", db.mgr.IdentityString(name), deleteDate, db.mgr.Placeholder(nil))
	args := []interface{}{utils.TimeToSQL(tm)}

	if c != nil {
		var where string
		var wargs []interface{}
		// MSSQL uses numbered placeholders (@pN); the SET clause already
		// consumed @p1 so WHERE placeholders must start at @p2. For "?"
		// placeholders the position is irrelevant.
		if db.dbtype == "sqlserver" {
			where, wargs = c.whereClause(db.mgr, &pos)
		} else {
			where, wargs = c.WhereClause(db.mgr)
		}
		if trimmed := strings.TrimSpace(where); trimmed != "" {
			s += fmt.Sprintf(" %s", trimmed)
			args = append(args, wargs...)
		}
	}
	return s, args
}

// RemoveMany removes all models of the specified type that match the criteria
// @param c
// @return int
// @return bool
func (db *DB) RemoveMany(m Modeller, c *Criteria, tx ...*sql.Tx) (int, error) {
	t := GetTableName(m)
	if !db.tableExists(t) {
		return 0, nil
	}
	r := db.Count(m, c)
	if r == 0 {
		return 0, nil
	}
	s := ""
	var args []interface{}
	if db.cfg.Deletable {
		s, args = db.massDelete(m, c)
	} else {
		s, args = db.massDisable(m, c)
	}
	err := db.executeQuery(s, args, tx...)
	return r, err
}

// columnType resolves the dialect-specific SQL column type for a field.
// SQL Server has no unsigned integer types, so unsigned fields are widened
// (INT -> BIGINT, BIGINT -> DECIMAL(20,0)) to preserve the full range.
func (db *DB) columnType(f field) string {
	switch db.cfg.Type {
	case "mysql", "mariadb":
		return mysqlFieldTypes[f.fType]
	case "sqlserver", "mssql":
		if f.unsigned {
			switch f.fType {
			case tInt:
				return "BIGINT"
			case tLong:
				return "DECIMAL(20,0)"
			}
		}
		return mssqlFieldTypes[f.fType]
	default:
		return sqliteFieldTypes[f.fType]
	}
}

// sized reports whether the field type accepts a "(size[,scale])" suffix in
// the current dialect. SQL Server only supports precision on DECIMAL and
// length on character types; a numeric "(n)" suffix is invalid T-SQL.
func (db *DB) sized(t fieldType) bool {
	switch db.cfg.Type {
	case "sqlserver", "mssql":
		return t == tDecimal || t == tString
	default:
		return true
	}
}

func (db *DB) tableDefinition(m Modeller) ([]string, bool) {
	sql := make([]string, 0, 3)

	n := GetTableName(m)
	db.mu.RLock()
	if _, ok := db.tableDef[n]; ok {
		db.mu.RUnlock()
		return nil, false
	}
	db.mu.RUnlock()

	db.mu.Lock()
	defer db.mu.Unlock()
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
	flds := getDefs(nm, true)
	db.tableDef[n] = flds
	if len(flds) == 0 {
		return nil, false
	}
	fldsStr := ""
	keys := make([]string, 0, 5)
	for _, f := range flds {
		if fldsStr != "" {
			fldsStr += ", "
		}
		ct := db.columnType(f)
		fldsStr += fmt.Sprintf("%s %s", db.mgr.IdentityString(f.name), ct)
		if f.size.Size > 0 && db.sized(f.fType) && !strings.Contains(ct, "(") {
			fldsStr += fmt.Sprintf("(%s)", f.size.String())
		}
		if f.fType == tString && f.size.Size == 0 {
			fldsStr += "(256)"
		}
		if f.unsigned && unsignedCapable(f.fType) &&
			(db.cfg.Type == "mysql" || db.cfg.Type == "mariadb") {
			fldsStr += " UNSIGNED"
		}
		if !f.allowNull {
			fldsStr += " NOT NULL"
		}
		if f.identity {
			fldsStr += " PRIMARY KEY"
		} else if f.key {
			keys = append(keys, f.name)
		}
	}
	sql = append(sql, fmt.Sprintf(db.mgr.TableCreate(), n, fldsStr))
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
	m, err := db.First(m, where.Equal("ID", *m.GetID()))
	return err
}

func Fetch[T Modeller](db *DB, criteria ...interface{}) ([]*T, error) {
	m := new(T)
	ms, err := db.Fetch(*m, criteria...)
	if err != nil {
		return nil, err
	}
	res := make([]*T, 0, len(ms))
	for _, m := range ms {
		res = append(res, utils.Ptr(m.(T)))
	}
	return res, nil
}

func First[T Modeller](db *DB, criteria ...interface{}) (*T, error) {
	r, err := Fetch[T](db, criteria...)
	if err != nil {
		return nil, err
	}
	if len(r) == 0 {
		return nil, NoResults("no results")
	}
	return r[0], err
}

func Range[T Modeller](db *DB, criteria ...interface{}) iter.Seq[*T] {
	m := new(T)
	return func(yield func(*T) bool) {
		for mdl := range db.Range(*m, criteria...) {
			if !yield(utils.Ptr(mdl.(T))) {
				return
			}
		}
	}
}

func FromID[T Modeller](db *DB, id string) (*T, error) {
	return First[T](db, where.Equal("ID", id))
}

func (db *DB) Connection() *sql.DB {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.db
}
