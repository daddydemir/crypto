package log

import "testing"

func TestInitLogger(t *testing.T) {
	InitLogger()

	LOG.Errorln("çalışmayabilir...")
}
