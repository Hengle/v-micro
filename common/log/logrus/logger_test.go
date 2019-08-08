package logrus

import (
	"testing"

	"github.com/fananchong/v-micro/common/log"
)

func TestLog(t *testing.T) {
	log := NewLogger(log.Name("testlog"))
	log.Infof("A")
	log.Info("B")
	log.Errorf("C")
	log.Error("D")
}
