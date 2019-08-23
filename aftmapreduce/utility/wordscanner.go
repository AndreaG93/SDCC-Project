package utility

import (
	"bufio"
	"strings"
)

func BuildWordScannerFromString(inputString string) *bufio.Scanner {

	output := bufio.NewScanner(strings.NewReader(inputString))
	output.Split(bufio.ScanWords)

	return output
}
