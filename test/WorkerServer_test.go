package test

import (
	"SDCC-Project/utility"
	"fmt"
	"os"
	"testing"
	"time"
)

type Logger struct {
	logFile *os.File
}

func NewLogger() *Logger {

	output := new(Logger)

	logFile, err := os.OpenFile("./log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	utility.CheckError(err)

	(*output).logFile = logFile
	return output
}

func (obj *Logger) log(message string) {

	message := fmt.Sprintf("TIME=\"%s\" -- ", time.Now().UTC().String())

	(*obj).logFile.WriteString(time.Now().UTC().String() + "sddsadasas\ndssadas\n\t\tsdasdasdas")
}

func Test_startWork(t *testing.T) {

	file, _ := os.OpenFile("./log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

}
