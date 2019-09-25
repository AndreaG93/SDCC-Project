package wordcount

import (
	"SDCC-Project/aftmapreduce/process"
	"SDCC-Project/aftmapreduce/utility"
	"fmt"
	"testing"
)

// Zookeeper must run on localhost!!
func Test_SaveDigestAssociation(t *testing.T) {

	utility.CheckError(process.Initialize(10, 0, process.WorkerProcessType, "localhost", []string{"127.0.0.1:2181"}))

	if err := SaveDigestAssociation("DIGEST002-1", "DIGEST"); err != nil {
		panic(err)
	}
	if err := SaveDigestAssociation("DIGEST001-1", "DIGEST"); err != nil {
		panic(err)
	}
	if err := SaveDigestAssociation("DIGEST001-3", "DIGEST"); err != nil {
		panic(err)
	}

	if output, err := GetDigestAssociationArray("DIGEST", 1); err != nil {
		panic(err)
	} else {
		fmt.Println(output)
	}

}
