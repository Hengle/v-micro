package logrus

import (
	"io"
	"os"
	"time"

	"github.com/fananchong/v-micro/common/log"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

type logger struct {
	*logrus.Logger
	opts log.Options
}

func (l *logger) Init(opt ...log.Option) (err error) {
	for _, o := range opt {
		o(&l.opts)
	}
	// 暂时关闭，目前只能修改源码才能达成，需要时再动手
	// https://github.com/sirupsen/logrus/issues/1004
	// l.SetReportCaller(true)
	l.SetLevel(logrus.Level(6 - int(l.opts.Level)))
	var output *rotatelogs.RotateLogs
	if output, err = rotatelogs.New(
		l.opts.Name+".%Y%m%d%H",
		rotatelogs.WithLinkName(l.opts.Name),
		rotatelogs.WithRotationTime(6*time.Hour),
	); err == nil {
		if l.opts.ToStdOut {
			l.SetOutput(io.MultiWriter(os.Stdout, output))
		} else {
			l.SetOutput(output)
		}
	}
	return
}

// String string
func (l *logger) String() string {
	return "logrus"
}

// NewLogger new
func NewLogger(opt ...log.Option) log.Logger {
	opts := log.Options{}
	for _, o := range opt {
		o(&opts)
	}

	l := &logger{logrus.New(), opts}
	return l
}
