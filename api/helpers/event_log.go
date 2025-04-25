package helpers

import (
	"github.com/semaphoreui/semaphore/db"
	"github.com/semaphoreui/semaphore/util"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type EventLogItem struct {
	IntegrationID int
	UserID        int
	ProjectID     int

	ObjectType  db.EventObjectType
	ObjectID    int
	Description string
}

type EventLogType string

const (
	EventLogCreate EventLogType = "create"
	EventLogUpdate EventLogType = "update"
	EventLogDelete EventLogType = "delete"
)

func EventLog(r *http.Request, action EventLogType, event EventLogItem) {
	record := db.Event{
		ObjectType:  &event.ObjectType,
		ObjectID:    &event.ObjectID,
		Description: &event.Description,
	}

	if event.IntegrationID > 0 {
		record.IntegrationID = &event.IntegrationID
	}

	if event.UserID > 0 {
		record.UserID = &event.UserID
	}

	if event.ProjectID > 0 {
		record.ProjectID = &event.ProjectID
	}

	logFields := record.ToFields()
	logFields["action"] = string(action)

	if _, err := Store(r).CreateEvent(record); err != nil {
		log.WithFields(logFields).Error("Failed to store event")
	}

	if err := util.Config.Log.Events.Write(logFields); err != nil {
		log.WithFields(logFields).Error("Failed to store event in log file")
	}
}
