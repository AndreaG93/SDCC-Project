package nodelogger

import (
	"SDCC-Project/aftmapreduce/node"
	"github.com/Sirupsen/logrus"
	"os"
)

type Logger struct {
	node string
	log  *logrus.Logger
}

func New() *Logger {

	output := new(Logger)
	(*output).log = logrus.New()

	file, err := os.OpenFile("./log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		(*output).log.Out = file
	} else {
		(*output).log.Info("Failed to log to file, using default stderr")
	}

	return output
}

func (obj *Logger) PrintMessage(message string) {
	(*obj).log.Infof("%s-%d :--> %s", node.GetStringProperty(node.TypeStringProperty), node.GetIntegerProperty(node.IDIntegerProperty), message)
}

func (obj *Logger) InvalidArgumentValue(argumentName string, argumentValue string) {
	(*obj).log.Infof("%s-%d :--> Invalid value for argument: %s: %v", node.GetStringProperty(node.TypeStringProperty), node.GetIntegerProperty(node.IDIntegerProperty), argumentName, argumentValue)
}
