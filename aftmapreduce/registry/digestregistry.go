package registry

import (
	"sync"
)

type DigestCache struct {
	table map[string][]string
}

var data *DigestCache
var once sync.Once

func GetDigestCacheInstance() *DigestCache {
	once.Do(func() {
		data = NewDigestCache()
	})
	return data
}

func NewDigestCache() *DigestCache {

	output := new(DigestCache)
	(*output).table = make(map[string][]string)

	return output
}

func (obj *DigestCache) Add(digest string, associatedDigest string) {

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

func (obj *DigestCache) GetAssociatedDigest(digest string) []string {
	return (*obj).table[digest]
}
