package stage_parsers

import (
	"github.com/semaphoreui/semaphore/db"
)

func MoveToNextStage(
	store db.Store,
	app db.TemplateApp,
	projectID int,
	currentState any,
	currentStage *db.TaskStage,
	currentOutput *db.TaskOutput,
	newOutput db.TaskOutput,
) (newStage *db.TaskStage, newState any, err error) {
	return
}
