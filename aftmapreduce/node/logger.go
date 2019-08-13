package node

import (
	"SDCC-Project/aftmapreduce/node/property"
	"github.com/Sirupsen/logrus"
	"os"
	"strings"
)

type Logger struct {
	log *logrus.Logger
}

func NewLogger() *Logger {

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

func (obj *Logger) PrintInfoTaskMessage(taskName string, taskMessage string) {

	nodeId := GetPropertyAsInteger(property.NodeID)
	nodeType := strings.ToUpper(GetPropertyAsString(property.NodeType))

	(*obj).log.Infof("NODE: \"%s-%d\" >> TASK: \"%s\" -->: %s", nodeType, nodeId, taskName, taskMessage)
}

func (obj *Logger) PrintErrorTaskMessage(taskName string, taskMessage string) {

	nodeId := GetPropertyAsInteger(property.NodeID)
	nodeType := strings.ToUpper(GetPropertyAsString(property.NodeType))

	(*obj).log.Errorf("NODE: \"%s-%d\" >> TASK: \"%s\" -->: %s", nodeType, nodeId, taskName, taskMessage)
}

func (obj *Logger) PrintPanicErrorTaskMessage(taskName string, taskMessage string) {

	nodeId := GetPropertyAsInteger(property.NodeID)
	nodeType := strings.ToUpper(GetPropertyAsString(property.NodeType))

	(*obj).log.Panicf("NODE: \"%s-%d\" >> TASK: \"%s\" -->: %s", nodeType, nodeId, taskName, taskMessage)
}

func (obj *Logger) PrintInfoStartingTaskMessage(taskName string) {
	(*obj).PrintInfoTaskMessage(taskName, "starting...")
}

func (obj *Logger) PrintInfoCompleteTaskMessage(taskName string) {
	(*obj).PrintInfoTaskMessage(taskName, "complete!")
}
