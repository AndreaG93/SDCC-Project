package file_split

import "testing"

func Test_SplitFile(t *testing.T) {

	if err := sss("./test.txt", 5); err != nil {
		panic(err)
	}
	/*
		if err := SplitFile("./test.txt", 5); err != nil {
			panic(err)
		}
	*/
}
