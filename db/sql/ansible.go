//go:build !pro

package sql

import (
	"errors"
	"github.com/semaphoreui/semaphore/db"
)

func (d *SqlDb) CreateAnsibleTaskHost(host db.AnsibleTaskHost) error {
	return nil
}

func (d *SqlDb) CreateAnsibleTaskError(error db.AnsibleTaskError) error {
	return nil
}

func (d *SqlDb) GetAnsibleTaskHosts(projectID int, taskID int) (res []db.AnsibleTaskHost, err error) {
	res = make([]db.AnsibleTaskHost, 0)
	return
}

func (d *SqlDb) GetAnsibleTaskErrors(projectID int, taskID int) (res []db.AnsibleTaskError, err error) {
	res = make([]db.AnsibleTaskError, 0)
	return
}
