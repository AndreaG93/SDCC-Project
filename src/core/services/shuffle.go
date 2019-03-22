package services

import "SDCC-Project-WorkerNode/src/core/data-structures"

func shuffle(input []*data_structures.WordTokenHashTable, output []*data_structures.WordTokenListArray) {

	size := uint(len(input))

	outputDataStructure := make([][]*data_structures.WordTokenList, size)

	for index := uint(0); index < size; index++ {
		outputDataStructure[index] = data_structures.BuildWordTokenListArray(size)
	}

	data_structures.WordTokenListArray{}.Built(size)

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
