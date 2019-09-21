package logger

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"os"
	"strings"
)

type Logger struct {
	logger               *logrus.Logger
	processID            int
	processType          string
	processWorkerGroupID int
}

func New(processID int, processType string, processWorkerGroupID int) (*Logger, error) {

	output := new(Logger)
	(*output).logger = logrus.New()
	(*output).processID = processID
	(*output).processType = strings.ToUpper(processType)
	(*output).processWorkerGroupID = processWorkerGroupID

	if file, err := os.OpenFile("./log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
		return nil, err
	} else {
		(*output).logger.Out = file
	}

	return output, nil
}

func (obj *Logger) PrintInfoLevelMessage(message string) {
	(*obj).logger.SetLevel(logrus.InfoLevel)
	(*obj).printMessage(message)
}

func (obj *Logger) PrintInfoLevelLabeledMessage(label string, message string) {
	(*obj).logger.SetLevel(logrus.InfoLevel)
	(*obj).printMessage(fmt.Sprintf("%s :: %s", label, message))
}

func (obj *Logger) PrintErrorLevelMessage(message string) {
	(*obj).logger.SetLevel(logrus.ErrorLevel)
	(*obj).printMessage(message)
}

func (obj *Logger) PrintPanicLevelMessage(message string) {
	(*obj).logger.SetLevel(logrus.PanicLevel)
	(*obj).printMessage(message)
}

func (obj *Logger) printMessage(message string) {

	if (*obj).processWorkerGroupID != -1 {
		(*obj).logger.Printf("%s-%d -- WPG: %d >> %s", (*obj).processType, (*obj).processID, (*obj).processWorkerGroupID, message)
	} else {
		(*obj).logger.Printf("%s-%d >> %s", (*obj).processType, (*obj).processID, message)
	}
}
