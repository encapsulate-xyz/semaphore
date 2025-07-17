//go:build pro

package sql

import (
	"github.com/Masterminds/squirrel"
	"github.com/semaphoreui/semaphore/db"
)

func (d *SqlDb) CreateAnsibleTaskHost(host db.AnsibleTaskHost) error {
	_, err := d.exec(
		"insert into task__ansible_host (project_id, task_id, host, failed, changed, ignored, ok, rescued, skipped, unreachable) "+
			"values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		host.ProjectID,
		host.TaskID,
		host.Host,
		host.Failed,
		host.Changed,
		host.Ignored,
		host.Ok,
		host.Rescued,
		host.Skipped,
		host.Unreachable,
	)
	if err != nil {
		return err
	}
	return nil
}

func (d *SqlDb) CreateAnsibleTaskError(error db.AnsibleTaskError) error {
	_, err := d.exec(
		"insert into task__ansible_error (project_id, task_id, host, task, error) "+
			"values (?, ?, ?, ?, ?)",
		error.ProjectID,
		error.TaskID,
		error.Host,
		error.Task,
		error.Error,
	)
	if err != nil {
		return err
	}
	return nil
}

func (d *SqlDb) GetAnsibleTaskHosts(projectID int, taskID int) (res []db.AnsibleTaskHost, err error) {
	q := squirrel.Select("*").
		From("task__ansible_host").
		Where("project_id=? and task_id=?", projectID, taskID)

	query, args, err := q.ToSql()

	if err != nil {
		return
	}

	_, err = d.selectAll(&res, query, args...)
	return
}

func (d *SqlDb) GetAnsibleTaskErrors(projectID int, taskID int) (res []db.AnsibleTaskError, err error) {
	q := squirrel.Select("*").
		From("task__ansible_error").
		Where("project_id=? and task_id=?", projectID, taskID)

	query, args, err := q.ToSql()

	if err != nil {
		return
	}

	_, err = d.selectAll(&res, query, args...)
	return
}
