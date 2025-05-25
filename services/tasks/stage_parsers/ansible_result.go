package stage_parsers

import (
	"github.com/semaphoreui/semaphore/db"
	"github.com/semaphoreui/semaphore/util"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
)

type AnsibleResultStageParser struct {
	state *AnsibleResultStageParserState
}

const ansibleResultMaker = "PLAY RECAP *****************************************"

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

	return strings.HasPrefix(output.Output, ansibleResultMaker)
}

func (p AnsibleResultStageParser) IsEnd(currentStage *db.TaskStage, output db.TaskOutput) bool {
	if currentStage == nil {
		return false
	}

	if currentStage.Type != db.TaskStagePrintResult {
		return false
	}

	return strings.TrimSpace(output.Output) == ""
}

//type AnsibleResultHost struct {
//	Host        string `json:"host"`
//	Ok          int    `json:"ok"`
//	Changed     int    `json:"changed"`
//	Unreachable int    `json:"unreachable"`
//	Failed      int    `json:"failed"`
//	Skipped     int    `json:"skipped"`
//	Rescued     int    `json:"rescued"`
//	Ignored     int    `json:"ignored"`
//}

var ansibleResultHostRE = regexp.MustCompile(
	`^([^\s]+)\s*:\s*` +
		`ok=(\d+)\s+` +
		`changed=(\d+)\s+` +
		`unreachable=(\d+)\s+` +
		`failed=(\d+)\s+` +
		`skipped=(\d+)\s+` +
		`rescued=(\d+)\s+` +
		`ignored=(\d+)$`)

func toInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
	}
	return v
}

type AnsibleResultStageParserState struct {
	//Hosts []AnsibleResultHost
}

func (p AnsibleResultStageParser) Parse(currentStage *db.TaskStage, output db.TaskOutput, store db.Store, projectID int) (ok bool, err error) {

	if currentStage == nil {
		return
	}

	if currentStage.Type != db.TaskStagePrintResult {
		return
	}

	ok = true

	line := util.ClearFromAnsiCodes(strings.TrimSpace(output.Output))

	if line == "" {
		return
	}

	if strings.HasPrefix(line, ansibleResultMaker) {
		return
	}

	m := ansibleResultHostRE.FindStringSubmatch(line)
	if m == nil {
		log.WithFields(log.Fields{
			"task_id": output.TaskID,
		}).Warnf("invalid ansible result host: %s", line)
		return
	}

	//p.state.Hosts = append(p.state.Hosts, AnsibleResultHost{
	//	Host:        m[1],
	//	Ok:          toInt(m[2]),
	//	Changed:     toInt(m[3]),
	//	Unreachable: toInt(m[4]),
	//	Failed:      toInt(m[5]),
	//	Skipped:     toInt(m[6]),
	//	Rescued:     toInt(m[7]),
	//	Ignored:     toInt(m[8]),
	//})

	err = store.CreateAnsibleTaskHost(db.AnsibleTaskHost{
		TaskID:      currentStage.TaskID,
		ProjectID:   projectID,
		Host:        m[1],
		Ok:          toInt(m[2]),
		Changed:     toInt(m[3]),
		Unreachable: toInt(m[4]),
		Failed:      toInt(m[5]),
		Skipped:     toInt(m[6]),
		Rescued:     toInt(m[7]),
		Ignored:     toInt(m[8]),
	})

	return
}

func (p AnsibleResultStageParser) Result() (res map[string]any) {
	res = make(map[string]any)
	//res["hosts"] = p.state.Hosts
	return
}

func (p AnsibleResultStageParser) State() any {
	return p.state
}
