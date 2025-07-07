package sql

import (
	"context"
	"database/sql"
	"github.com/go-gorp/gorp/v3"
	"github.com/semaphoreui/semaphore/util"
	"sync"
)

type FakeLocker struct {
}

func (f FakeLocker) Lock() {
}

func (f FakeLocker) Unlock() {
}

type CommonSql struct {
	db    *gorp.DbMap
	mutex sync.Locker
}

func Create(dialectStr string, sqlDb *sql.DB) *CommonSql {
	var dialect gorp.Dialect

	switch dialectStr {
	case util.DbDriverMySQL:
		dialect = gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"}
	case util.DbDriverPostgres:
		dialect = gorp.PostgresDialect{}
	case util.DbDriverSQLite:
		dialect = gorp.SqliteDialect{}
	}

	var mutex sync.Locker
	if dialectStr == util.DbDriverSQLite {
		mutex = &sync.Mutex{}
	} else {
		mutex = &FakeLocker{}
	}

	return &CommonSql{db: &gorp.DbMap{Db: sqlDb, Dialect: dialect}, mutex: mutex}
}

func (c CommonSql) Lock() {
	c.mutex.Lock()
}

func (c CommonSql) Unlock() {
	c.mutex.Unlock()
}

func (c CommonSql) QueryRow(query string, args ...interface{}) *sql.Row {
	c.Lock()
	defer c.Unlock()
	return c.db.QueryRow(query, args...)
}

func (c CommonSql) Insert(list ...interface{}) error {
	c.Lock()
	defer c.Unlock()

	return c.db.Insert(list...)
}

func (c CommonSql) SelectOne(holder interface{}, query string, args ...interface{}) error {
	c.Lock()
	defer c.Unlock()

	return c.db.SelectOne(holder, query, args...)
}

func (c CommonSql) PrepareQuery(query string) string {
	return query
}

func (c CommonSql) Exec(query string, args ...interface{}) (sql.Result, error) {
	c.Lock()
	defer c.Unlock()

	return c.db.Exec(query, args...)
}

func (c CommonSql) Select(i interface{}, query string, args ...interface{}) ([]interface{}, error) {
	c.Lock()
	defer c.Unlock()

	return c.db.Select(i, query, args...)
}

func (c CommonSql) SelectInt(query string, args ...interface{}) (int64, error) {
	c.Lock()
	defer c.Unlock()

	return c.db.SelectInt(query, args...)
}

func (c CommonSql) AddTableWithName(i interface{}, name string) *gorp.TableMap {
	c.Lock()
	defer c.Unlock()

	return c.db.AddTableWithName(i, name)
}

func (c CommonSql) Dialect() gorp.Dialect {
	return c.db.Dialect
}

func (c CommonSql) Close() error {
	c.Lock()
	defer c.Unlock()

	return c.db.Db.Close()
}

func (c CommonSql) Begin() (*gorp.Transaction, error) {
	return c.db.Begin()
}

func (c CommonSql) WithContext(ctx context.Context) gorp.SqlExecutor {
	c.Lock()
	defer c.Unlock()

	return c.db.WithContext(ctx)
}

func (c CommonSql) Get(i interface{}, keys ...interface{}) (interface{}, error) {
	c.Lock()
	defer c.Unlock()

	return c.db.Get(i, keys...)
}

func (c CommonSql) Update(list ...interface{}) (int64, error) {
	c.Lock()
	defer c.Unlock()

	return c.db.Update(list...)
}

func (c CommonSql) Delete(list ...interface{}) (int64, error) {
	c.Lock()
	defer c.Unlock()

	return c.db.Delete(list...)
}

func (c CommonSql) SelectNullInt(query string, args ...interface{}) (sql.NullInt64, error) {
	c.Lock()
	defer c.Unlock()

	return c.db.SelectNullInt(query, args...)
}

func (c CommonSql) SelectFloat(query string, args ...interface{}) (float64, error) {
	c.Lock()
	defer c.Unlock()

	return c.db.SelectFloat(query, args...)
}

func (c CommonSql) SelectNullFloat(query string, args ...interface{}) (sql.NullFloat64, error) {
	c.Lock()
	defer c.Unlock()

	return c.db.SelectNullFloat(query, args...)
}

func (c CommonSql) SelectStr(query string, args ...interface{}) (string, error) {
	c.Lock()
	defer c.Unlock()

	return c.db.SelectStr(query, args...)
}

func (c CommonSql) SelectNullStr(query string, args ...interface{}) (sql.NullString, error) {
	c.Lock()
	defer c.Unlock()

	return c.db.SelectNullStr(query, args...)
}

func (c CommonSql) Query(query string, args ...interface{}) (*sql.Rows, error) {
	c.Lock()
	defer c.Unlock()

	return c.db.Query(query, args...)
}
