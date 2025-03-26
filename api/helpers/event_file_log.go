//go:build !pro

package helpers

import (
	"github.com/semaphoreui/semaphore/db"
)

// AppendEventToLog opens (or creates) the log file in append mode and writes the Event in key=value format.
func appendEventToLog(event db.Event) error {
	return nil
}
