package stage_parsers

import (
	"github.com/semaphoreui/semaphore/db"
	"strings"
)

type AnsibleRunningStageParserFailedTask struct {
	Task   string
	Host   string
	Answer string
}

type AnsibleRunningStageParserState struct {
	CurrentTask       string
	CurrentFailedHost string
	CurrentHostAnswer string
	FailedTasks       []AnsibleRunningStageParserFailedTask
}

type AnsibleRunningStageParser struct {
	state *AnsibleRunningStageParserState
}

func (p AnsibleRunningStageParser) State() any {
	return p.state
}

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

const ansibleTaskMaker = "TASK ["
const failedTaskMaker = "fatal: ["

func (p AnsibleRunningStageParser) Parse(output db.TaskOutput) error {

	if strings.HasPrefix(output.Output, ansibleTaskMaker) {
		if p.state.CurrentFailedHost != "" {
			p.state.FailedTasks = append(p.state.FailedTasks, AnsibleRunningStageParserFailedTask{
				Task:   p.state.CurrentTask,
				Host:   p.state.CurrentFailedHost,
				Answer: p.state.CurrentHostAnswer,
			})
		}

		end := strings.Index(output.Output, "]")
		p.state.CurrentTask = output.Output[len(ansibleTaskMaker):end]
		p.state.CurrentFailedHost = ""
		p.state.CurrentHostAnswer = ""
	} else if strings.HasPrefix(output.Output, failedTaskMaker) {
		end := strings.Index(output.Output, "]")
		p.state.CurrentFailedHost = output.Output[len(ansibleTaskMaker):end]
		p.state.CurrentHostAnswer = ""
	} else if p.state.CurrentFailedHost != "" {
		if output.Output == "" {
			if p.state.CurrentFailedHost != "" {
				p.state.FailedTasks = append(p.state.FailedTasks, AnsibleRunningStageParserFailedTask{
					Task:   p.state.CurrentTask,
					Host:   p.state.CurrentFailedHost,
					Answer: p.state.CurrentHostAnswer,
				})
			}
			p.state.CurrentFailedHost = ""
			p.state.CurrentHostAnswer = ""
		} else {
			p.state.CurrentHostAnswer += "\n" + output.Output
		}
	}

	// Implement the parsing logic for Ansible results
	return nil
}

func (p AnsibleRunningStageParser) Result() (res map[string]any) {

	return nil
}
