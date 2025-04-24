package stage_parsers

import (
	"github.com/semaphoreui/semaphore/db"
	"strings"
)

type AnsibleResultStageParser struct{}

func (p AnsibleResultStageParser) NeedParse() bool {
	return true
}

func (p AnsibleResultStageParser) IsStart(currentStage *db.TaskStage, output db.TaskOutput) bool {
	if currentStage == nil {
		return false
	}

	if currentStage.Type != db.TaskStageRunning {
		return false
	}

	return strings.HasPrefix(output.Output, "PLAY RECAP *****************************************")
}

func (p AnsibleResultStageParser) IsEnd(currentStage *db.TaskStage, output db.TaskOutput) bool {
	return false
}

func (p AnsibleResultStageParser) Parse(outputs []db.TaskOutput) (map[string]any, error) {
	// Implement the parsing logic for Ansible results
	return nil, nil
}

type AnsibleRunningStageParser struct{}

func (p AnsibleRunningStageParser) NeedParse() bool {
	return false
}

func (p AnsibleRunningStageParser) IsStart(currentStage *db.TaskStage, output db.TaskOutput) bool {

	if currentStage == nil {
		return false
	}

	if currentStage.Type != db.TaskStageInit {
		return false
	}

	return strings.HasPrefix(output.Output, "PLAY [")
}

func (p AnsibleRunningStageParser) IsEnd(currentStage *db.TaskStage, output db.TaskOutput) bool {
	return false
}

func (p AnsibleRunningStageParser) Parse(outputs []db.TaskOutput) (map[string]any, error) {
	// Implement the parsing logic for Ansible results
	return nil, nil
}
