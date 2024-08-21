package test

import (
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	file, err := os.OpenFile("test.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	logger := logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.SetOutput(file)

	logger.WithField("keybaru", "val").Info("info")
	logger.Warning("warn")
	logger.Error("err")
}
