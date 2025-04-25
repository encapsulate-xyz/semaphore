//go:build !pro

package util

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func appendToFileLog(row log.Fields, logger *lumberjack.Logger) error {
	return nil
}
