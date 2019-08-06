package logrus

import (
	"testing"
)

func TestLog(t *testing.T) {
	log := NewLogger("out")
	log.Infof("A")
	log.Infoln("B")
	log.Errorf("C")
	log.Errorln("D")
}
