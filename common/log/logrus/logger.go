package logrus

import (
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

// NewLogger new
func NewLogger(logName string) *logrus.Logger {
	l := logrus.New()
	l.SetReportCaller(true)
	if output, err := rotatelogs.New(
		logName+".%Y%m%d%H",
		rotatelogs.WithLinkName(logName),
		rotatelogs.WithRotationTime(6*time.Hour),
	); err == nil {
		l.SetOutput(output)
	}
	return l
}
