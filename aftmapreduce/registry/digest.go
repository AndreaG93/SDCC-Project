package registry

type DigestRegistry struct {
	table map[string][]string
}

func NewDigestRegistry() *DigestRegistry {

	output := new(DigestRegistry)
	(*output).table = make(map[string][]string)

	return output
}

func (obj *DigestRegistry) Add(digest string, associatedDigest string) {

	associatedDigests := (*obj).table[digest]

	if associatedDigests == nil {

		associatedDigests = make([]string, 1)
		associatedDigests[0] = associatedDigest

		(*obj).table[digest] = associatedDigests
		return
	} else {

		found := false

		for _, elem := range associatedDigests {
			if elem == digest {
				found = true
			}
		}

		if !found {
			(*obj).table[digest] = append(associatedDigests, associatedDigest)
		}
	}
}

func (obj *DigestRegistry) GetAssociatedDigest(digest string) []string {
	return (*obj).table[digest]
}
