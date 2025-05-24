package stage_parsers

import (
	"github.com/semaphoreui/semaphore/db"
	"strings"
)

type StageResultParser interface {
	IsStart(currentStage *db.TaskStage, output db.TaskOutput) bool
	IsEnd(currentStage *db.TaskStage, output db.TaskOutput) bool
	NeedParse() bool
	Parse(output db.TaskOutput) error
	State() any
	Result() map[string]any
}

func GetStageResultParser(app db.TemplateApp, stageType db.TaskStageType, state any) StageResultParser {

	if stageType == db.TaskStageInit {
		return &InitStageParser{}
	}

	switch app {
	case db.AppAnsible:
		switch stageType {
		case db.TaskStageRunning:

			if state == nil {
				state = &AnsibleRunningStageParserState{}
			}

			return &AnsibleRunningStageParser{
				state: state.(*AnsibleRunningStageParserState),
			}
		case db.TaskStagePrintResult:

			if state == nil {
				state = &AnsibleResultStageParserState{}
			}

			return &AnsibleResultStageParser{
				state: state.(*AnsibleResultStageParserState),
			}
		}
	}

	return nil
}

func GetAllTaskStages(app db.TemplateApp) []db.TaskStageType {
	switch app {
	case db.AppAnsible:
		return []db.TaskStageType{
			db.TaskStageInit,
			db.TaskStageRunning,
			db.TaskStagePrintResult,
		}
	}

	return nil
}

type InitStageParser struct{}

func (p InitStageParser) State() any {
	return nil
}

func (p InitStageParser) NeedParse() bool {
	return false
}

func (p InitStageParser) IsStart(currentStage *db.TaskStage, output db.TaskOutput) bool {
	if currentStage != nil {
		return false
	}

	return strings.HasPrefix(output.Output, "Run TaskRunner with template: ")
}

func (p InitStageParser) IsEnd(currentStage *db.TaskStage, output db.TaskOutput) bool {
	return false
}

func (p InitStageParser) Parse(output db.TaskOutput) error {
	return nil
}

func (p InitStageParser) Result() map[string]any {
	return nil
}
