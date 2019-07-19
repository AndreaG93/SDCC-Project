package Data

type RawInput interface {
	Split() []Split
	MapOutputRawDataToReduceInputData(mapOutputRawData [][]byte) []Split
	ReduceOutputRawDataToFinalOutput(ReduceOutputRawData [][]byte) []Split
}
