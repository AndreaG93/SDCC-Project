package services

import "SDCC-Project-WorkerNode/src/core/data-structuresttt"

func shuffle(input []*data_structuresttt.WordTokenHashTable, output []*data_structuresttt.WordTokenListArray) {

	size := uint(len(input))

	outputDataStructure := make([][]*data_structuresttt.WordTokenList, size)

	for index := uint(0); index < size; index++ {
		outputDataStructure[index] = data_structuresttt.BuildWordTokenListArray(size)
	}

	data_structuresttt.WordTokenListArray{}.Built(size)

	/*
		for x := uint(0); x < size; x++ {
			for y :=  uint(0); y < size; y++ {

				currentWordTokenList := input[x].GetWordTokenListAt(y)
				(outputDataStructure[x])[y] = currentWordTokenList

			}


		}

		for x := 0; x < core.AvailableWorkersNumber; x++ {
			for y := 0; y < core.AvailableWorkersNumber; y++ {
				pOutput[y].Data[x] = pInput[x].Data[y]
			}
		}
	*/
}
