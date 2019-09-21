package wordcount

import (
	"SDCC-Project/aftmapreduce/process"
	"fmt"
	"math"
	"strings"
)

func getSplits(guid string, splitsAmount int) ([]string, error) {

	if text, err := downloadTextFromCloud(guid); err != nil {
		return nil, err
	} else {
		process.GetLogger().PrintInfoLevelMessage(fmt.Sprintf("Start SPLIT phase :: splitsAmount %d", splitsAmount))
		return divideIntoSplits(text, splitsAmount), nil
	}
}

func downloadTextFromCloud(guid string) (string, error) {

	if rawData, err := (*process.GetStorageKeyValueRegister()).Get(guid); err == nil {
		return string(rawData), nil
	} else {
		return "", err
	}
}

func divideIntoSplits(input string, splitsAmount int) []string {

	output := make([]string, splitsAmount)
	outputIndex := 0
	splitSize := int(math.Floor(float64(len(input) / splitsAmount)))

	currentSplitLowerLimit := 0
	currentSplitUpperLimit := splitSize

	for {
		currentChar := string(input[currentSplitUpperLimit])

		if strings.Compare(currentChar, " ") == 0 {

			output[outputIndex] = input[currentSplitLowerLimit:currentSplitUpperLimit]
			outputIndex++

			if currentSplitUpperLimit+splitSize >= len(input) {
				output[outputIndex] = input[currentSplitUpperLimit:]
				break
			} else {
				currentSplitLowerLimit = currentSplitUpperLimit
				currentSplitUpperLimit += splitSize
			}

		} else {

			if currentSplitUpperLimit+1 == len(input) {
				output[outputIndex] = input[currentSplitLowerLimit:]
				break
			} else {
				currentSplitUpperLimit++
			}
		}
	}

	return output
}
