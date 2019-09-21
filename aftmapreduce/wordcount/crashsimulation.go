package wordcount

import (
	"SDCC-Project/aftmapreduce/process"
	"SDCC-Project/aftmapreduce/process/property"
)

func crash() {
	if process.GetPropertyAsInteger(property.NodeID) == 1 {
		panic("Simulated Crash")
	}

	if process.GetPropertyAsInteger(property.NodeID) == 0 {
		panic("Simulated Crash")
	}
}

func isOccurredAnArbitraryCrash() bool {
	if process.GetPropertyAsInteger(property.NodeID) == 4 {
		return true
	} else {
		return false
	}
}
