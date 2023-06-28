package log

import (
	"github.com/daddydemir/crypto/config"
	"os"
	"time"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var LOG *logrus.Logger
var LogFile *os.File

func InitLogger() {

	logPath := config.Get("LOG_PATH")
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	err := error(nil)
	LogFile, err = os.OpenFile(getLogFilePath(logPath), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("log dosyası açılamadı:", err)
	}

	log.AddHook(lfshook.NewHook(
		lfshook.PathMap{
			logrus.InfoLevel:  getLogFilePath(logPath + "info/"),
			logrus.ErrorLevel: getLogFilePath(logPath + "error/"),
		},
		&logrus.JSONFormatter{},
	))

	log.SetOutput(LogFile)

	LOG = log

	LOG.Infoln("log test")
}

func getLogFilePath(basePath string) string {
	currentDate := time.Now().Format("2006-01-02")
	return basePath + currentDate + ".log"
}

func Errorln(args ...interface{}) {
	LOG.Logln(2, args...)
}

func Infoln(args ...interface{}) {
	LOG.Logln(4, args...)
}

func Fatal(args ...interface{}) {
	LOG.Log(1, args...)
	LOG.Exit(1)
}
