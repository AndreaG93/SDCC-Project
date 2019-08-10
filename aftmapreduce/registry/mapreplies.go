package registry

type sameDigestReply struct {
	sameDigestReplyNodeIds []int
	MappedDataSizes        map[int]int
}

type MapReplies struct {
	replies                      map[string]*sameDigestReply
	requiredNumberOfMatches      int
	mostMatchedWorkerReplyDigest string
}

func NewMapReply(requiredNumberOfMatches int) *MapReplies {

	output := new(MapReplies)

	(*output).replies = make(map[string]*sameDigestReply)
	(*output).requiredNumberOfMatches = requiredNumberOfMatches

	return output
}

func (obj *MapReplies) Add(digest string, nodeId int, mappedDataSizes map[int]int) bool {

	reply := (*obj).replies[digest]

	if reply == nil {

		reply = new(sameDigestReply)

		(*reply).sameDigestReplyNodeIds = make([]int, 0)
		(*reply).MappedDataSizes = mappedDataSizes

		(*obj).replies[digest] = reply
	}

	(*reply).sameDigestReplyNodeIds = append((*reply).sameDigestReplyNodeIds, nodeId)

	if len((*reply).sameDigestReplyNodeIds) == (*obj).requiredNumberOfMatches {
		(*obj).mostMatchedWorkerReplyDigest = digest
		return true
	}

	return false
}

func (obj *MapReplies) GetMostMatchedReply() (string, []int, map[int]int) {

	mostMatchedReply := (*obj).replies[(*obj).mostMatchedWorkerReplyDigest]

	return (*obj).mostMatchedWorkerReplyDigest, mostMatchedReply.sameDigestReplyNodeIds, mostMatchedReply.MappedDataSizes
}
