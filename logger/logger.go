package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// L is the global logger instance.
var L *logrus.Logger = logrus.StandardLogger()

// famous logger key
const (
	App       = "appid"
	Sub       = "sub"
	UID       = "uid"
	Email     = "email"
	Phone     = "phone"
	Groups    = "groups"
)

// Production initialize the global logger for production purpose
func Production(level, writer string) error {
	var (
		w   io.Writer
		err error
	)

	if writer == "stdout" {
		w = os.Stdout
	} else {
		w, err = os.OpenFile(writer, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			fmt.Printf("failed to open log file %s: %v\n", writer, err)
			os.Exit(1)
		}
	}

	L = logrus.StandardLogger()
	L.Out = w

	l, err := logrus.ParseLevel(level)
	if err != nil {
		l = logrus.InfoLevel
	}
	L.SetLevel(l)

	return nil
}
