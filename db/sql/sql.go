package sql

import (
	"github.com/go-gorp/gorp/v3"
)

type Sql interface {
	gorp.SqlExecutor

	AddTableWithName(i interface{}, name string) *gorp.TableMap

	Dialect() gorp.Dialect

	Close() error

	Begin() (*gorp.Transaction, error)

	Lock()

	Unlock()
}
