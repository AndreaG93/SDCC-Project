package nodelogger

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"os"
)

type Logger struct {
	node string
	log  *logrus.Logger
}

func New(nodeID int, nodeType string) *Logger {
	output := new(Logger)
	(*output).log = logrus.New()
	(*output).node = fmt.Sprintf("%s %d", nodeType, nodeID)

	file, err := os.OpenFile(fmt.Sprintf("%s-%d.log", nodeType, nodeID), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		(*output).log.Out = file
	} else {
		(*output).log.Info("Failed to log to file, using default stderr")
	}

	return output
}

func (obj *Logger) PrintMessage(message string) {
	(*obj).log.Infof("%s :--> %s", (*obj).node, message)
}

func (obj *Logger) InvalidArgumentValue(argumentName string, argumentValue string) {
	(*obj).log.Infof("%s :--> Invalid value for argument: %s: %v", (*obj).node, argumentName, argumentValue)
}
