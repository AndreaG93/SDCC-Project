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
	myMessage := (*obj).getFormattedMessage(message)
	(*obj).logger.Info(myMessage)
}

func (obj *Logger) PrintInfoLevelLabeledMessage(label string, message string) {
	myMessage := (*obj).getFormattedMessage(fmt.Sprintf("%s :: %s", label, message))
	(*obj).logger.Info(myMessage)
}

func (obj *Logger) PrintErrorLevelMessage(message string) {
	myMessage := (*obj).getFormattedMessage(message)
	(*obj).logger.Error(myMessage)
}

func (obj *Logger) PrintPanicLevelMessage(message string) {
	myMessage := (*obj).getFormattedMessage(message)
	(*obj).logger.Panic(myMessage)
}

func (obj *Logger) getFormattedMessage(message string) string {

	if (*obj).processWorkerGroupID != -1 {
		return fmt.Sprintf("%s-%d -- WPG: %d >> %s", (*obj).processType, (*obj).processID, (*obj).processWorkerGroupID, message)
	} else {
		return fmt.Sprintf("%s-%d >> %s", (*obj).processType, (*obj).processID, message)
	}
}
