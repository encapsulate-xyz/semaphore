package factory

import (
	"github.com/semaphoreui/semaphore/db"
	"github.com/semaphoreui/semaphore/pro/db/sql"
)

func NewTerraformStore(store db.Store) db.TerraformStore {
	return &sql.TerraformStoreImpl{}
}
