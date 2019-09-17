package wordcount

import (
	"fmt"
	"testing"
)

func Test_Splits(t *testing.T) {

	output := divideIntoSplits("Andrea Andrea Andrea Andrea Andrea Andrea", 5)
	for _, x := range output {
		fmt.Println(x)
	}
}
