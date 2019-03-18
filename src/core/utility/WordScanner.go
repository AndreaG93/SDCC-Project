package utility

import (
	"bufio"
	"os"
	"strings"
)

func BuildWordScannerFromAnOpenedFile(inputFile *os.File) *bufio.Scanner {

	output := bufio.NewScanner(inputFile)
	output.Split(bufio.ScanWords)

	return output
}

func BuildWordScannerFromString(inputString string) *bufio.Scanner {

	output := bufio.NewScanner(strings.NewReader(inputString))
	output.Split(bufio.ScanWords)

	return output
}
