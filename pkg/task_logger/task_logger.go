package task_logger

import (
	"os/exec"
	"strings"
	"time"
)

type TaskStatus string

const (
	TaskWaitingStatus       TaskStatus = "waiting"
	TaskStartingStatus      TaskStatus = "starting"
	TaskWaitingConfirmation TaskStatus = "waiting_confirmation"
	TaskConfirmed           TaskStatus = "confirmed"
	TaskRejected            TaskStatus = "rejected"
	TaskRunningStatus       TaskStatus = "running"
	TaskStoppingStatus      TaskStatus = "stopping"
	TaskStoppedStatus       TaskStatus = "stopped"
	TaskSuccessStatus       TaskStatus = "success"
	TaskFailStatus          TaskStatus = "error"
)

func (s TaskStatus) String() string {
	switch s {
	case TaskWaitingStatus:
		return "waiting"
	case TaskStartingStatus:
		return "starting"
	case TaskWaitingConfirmation:
		return "waiting_confirmation"
	case TaskConfirmed:
		return "confirmed"
	case TaskRejected:
		return "rejected"
	case TaskRunningStatus:
		return "running"
	case TaskStoppingStatus:
		return "stopping"
	case TaskStoppedStatus:
		return "stopped"
	case TaskSuccessStatus:
		return "success"
	case TaskFailStatus:
		return "error"
	default:
		return "unknown"
	}
}

func (s TaskStatus) IsValid() bool {
	switch s {
	case TaskWaitingStatus,
		TaskStartingStatus,
		TaskWaitingConfirmation,
		TaskConfirmed,
		TaskRejected,
		TaskRunningStatus,
		TaskStoppingStatus,
		TaskStoppedStatus,
		TaskSuccessStatus,
		TaskFailStatus:
		return true
	}
	return false
}

func (s TaskStatus) IsNotifiable() bool {
	return s == TaskSuccessStatus || s == TaskFailStatus || s == TaskWaitingConfirmation
}

func (s TaskStatus) Format() (res string) {

	switch s {
	case TaskFailStatus:
		res += "❌"
	case TaskSuccessStatus:
		res += "✅"
	case TaskStoppedStatus:
		res += "⏹️"
	case TaskWaitingConfirmation:
		res += "⚠️"
	default:
		res += "❓"
	}

	// to avoid email content injection issue
	res += strings.ToUpper(s.String())

	return
}

func (s TaskStatus) IsFinished() bool {
	return s == TaskStoppedStatus || s == TaskSuccessStatus || s == TaskFailStatus
}

type StatusListener func(status TaskStatus)
type LogListener func(new time.Time, msg string)

type Logger interface {
	Log(msg string)
	Logf(format string, a ...any)
	LogWithTime(now time.Time, msg string)
	LogfWithTime(now time.Time, format string, a ...any)
	LogCmd(cmd *exec.Cmd)
	SetStatus(status TaskStatus)
	AddStatusListener(l StatusListener)
	AddLogListener(l LogListener)

	SetCommit(hash, message string)

	WaitLog()
}
