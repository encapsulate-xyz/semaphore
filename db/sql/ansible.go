//go:build !pro

package sql

import (
	"errors"
	"github.com/semaphoreui/semaphore/db"
)

func (d *SqlDb) CreateAnsibleTaskHost(host db.AnsibleTaskHost) error {
	return errors.New("not implemented")
}

func (d *SqlDb) CreateAnsibleTaskError(error db.AnsibleTaskError) error {
	return errors.New("not implemented")
}

func (d *SqlDb) GetAnsibleTaskHosts(projectID int, taskID int) (res []db.AnsibleTaskHost, err error) {
	err = errors.New("not implemented")
	return
}

func (d *SqlDb) GetAnsibleTaskErrors(projectID int, taskID int) (res []db.AnsibleTaskError, err error) {
	err = errors.New("not implemented")
	return
}
