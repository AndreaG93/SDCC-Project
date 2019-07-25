package data

type TransientData interface {
	PerformTask() (string, []byte, error)
}

type ClientData interface {
	GetDigest() string
	ToByte() []byte
	GetTypeName() string
	Split() ([]TransientData, error)
	Shuffle(rawDataFromMapTask [][]byte) []TransientData
	CollectResults(rawDataFromReduceTask [][]byte) string
}
