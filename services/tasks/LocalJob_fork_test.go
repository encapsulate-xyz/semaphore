package tasks

import (
	"encoding/json"
	"os/exec"
	"testing"
	"time"

	"github.com/semaphoreui/semaphore/db"
	"github.com/semaphoreui/semaphore/db_lib"
	"github.com/semaphoreui/semaphore/pkg/task_logger"
	"github.com/semaphoreui/semaphore/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type forkTestLogger struct{}

func (forkTestLogger) Log(string)                                   {}
func (forkTestLogger) Logf(string, ...any)                          {}
func (forkTestLogger) LogWithTime(time.Time, string)                {}
func (forkTestLogger) LogfWithTime(time.Time, string, ...any)       {}
func (forkTestLogger) LogCmd(*exec.Cmd)                             {}
func (forkTestLogger) SetStatus(task_logger.TaskStatus)             {}
func (forkTestLogger) AddStatusListener(task_logger.StatusListener) {}
func (forkTestLogger) AddLogListener(task_logger.LogListener)       {}
func (forkTestLogger) SetCommit(string, string)                     {}
func (forkTestLogger) WaitLog()                                     {}

// c53635e4: a task only runs from main, whichever of the repository, the template or the task names the branch.
func TestLocalJob_prepareRunOnlyMain(t *testing.T) {
	original := util.Config
	t.Cleanup(func() { util.Config = original })
	util.Config = &util.ConfigType{TmpPath: t.TempDir()}

	dev := "dev"
	cases := map[string]LocalJob{
		"repository on dev": {Repository: db.Repository{GitBranch: "dev"}},
		"template sets dev": {Repository: db.Repository{GitBranch: "main"}, Template: db.Template{GitBranch: &dev}},
		"task sets dev":     {Repository: db.Repository{GitBranch: "main"}, Task: db.Task{GitBranch: &dev}},
	}
	for name, job := range cases {
		t.Run(name, func(t *testing.T) {
			job.Logger = forkTestLogger{}
			err := job.prepareRun(db_lib.LocalAppInstallingArgs{})
			require.Error(t, err)
			assert.Contains(t, err.Error(), "branch must be 'main'")
		})
	}
}

// 95898999: a task's secret variables reach the playbook only when the template declares them as survey variables.
func TestLocalJob_secretVarsOnlySurveyVars(t *testing.T) {
	original := util.Config
	t.Cleanup(func() { util.Config = original })
	util.Config = &util.ConfigType{}

	job := LocalJob{
		Logger:      forkTestLogger{},
		Template:    db.Template{SurveyVars: []db.SurveyVar{{Name: "token"}}},
		Environment: db.Environment{JSON: `{"env": "prod"}`},
		Secret:      `{"token": "s3cret", "ansible_become_password": "injected"}`,
	}
	str, err := job.getEnvironmentExtraVarsJSON("someone", nil)
	require.NoError(t, err)

	vars := map[string]any{}
	require.NoError(t, json.Unmarshal([]byte(str), &vars))
	assert.Equal(t, "prod", vars["env"])
	assert.Equal(t, "s3cret", vars["token"])
	assert.NotContains(t, vars, "ansible_become_password")
	assert.Equal(t, "{}", job.Secret, "the secret is cleared once read")
}
