package Input

type MiddleInput interface {
	PerformTask() (string, []byte, error)
}

type ApplicationInput interface {
	Split() ([]MiddleInput, error)
	Shuffle(rawDataFromMapTask [][]byte) []MiddleInput
	CollectResults(rawDataFromReduceTask [][]byte) string
}
