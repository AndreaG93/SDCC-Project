package wordcount

import (
	"SDCC-Project/aftmapreduce/process"
	"SDCC-Project/aftmapreduce/utility"
	"fmt"
)

type Receive struct {
}

type ReceiveInput struct {
	Data                 []byte
	ReceivedDataDigest   string
	AssociatedDataDigest string
}

type ReceiveOutput struct {
}

const (
	digestAssociationArrayLabel = "DIGEST-ASSOCIATIONS"
)

func (x *Receive) Execute(input ReceiveInput, output *ReceiveOutput) error {

	process.GetLogger().PrintInfoLevelLabeledMessage(ReceiveTaskName, fmt.Sprintf("Received Data Digest: %s Associated to Data Digest: %s", input.ReceivedDataDigest, input.AssociatedDataDigest))

	if err := process.GetDataRegistry().Set(input.ReceivedDataDigest, input.Data); err != nil {
		return err
	}

	if input.AssociatedDataDigest != "" {
		if err := SaveDigestAssociation(input.ReceivedDataDigest, input.AssociatedDataDigest); err != nil {
			return err
		}
	}

	return nil
}

func GetDigestAssociationArray(localDigest string) ([]string, error) {

	var output []string

	key := fmt.Sprintf("%s-%s", localDigest, digestAssociationArrayLabel)
	rawData := process.GetDataRegistry().Get(key)

	if err := utility.Decoding(rawData, &output); err != nil {
		return nil, err
	} else {
		return output, nil
	}
}

func SaveDigestAssociation(digest string, localDigest string) error {

	var digestAssociationArray []string

	key := fmt.Sprintf("%s-%s", localDigest, digestAssociationArrayLabel)
	rawData := process.GetDataRegistry().Get(key)

	if rawData == nil {

		digestAssociationArray := make([]string, 1)
		digestAssociationArray[0] = digest

	} else {

		if err := utility.Decoding(rawData, &digestAssociationArray); err != nil {
			return err
		} else {

			for _, elem := range digestAssociationArray {
				if elem == digest {
					return nil
				}
			}

			digestAssociationArray = append(digestAssociationArray, digest)
		}
	}

	if rawData, err := utility.Encoding(digestAssociationArray); err != nil {
		return err
	} else {
		return process.GetDataRegistry().Set(key, rawData)
	}
}
