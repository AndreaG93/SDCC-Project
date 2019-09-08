package reply

type sameDigestMapReply struct {
	sameDigestReplyNodeIds []int
	dataSize               map[int]int
}

type MapReplyRegistry struct {
	registry                     map[string]*sameDigestMapReply
	requiredNumberOfMatches      int
	mostMatchedWorkerReplyDigest string
}

func NewMapReplyRegistry(requiredNumberOfMatches int) *MapReplyRegistry {

	output := new(MapReplyRegistry)

	(*output).registry = make(map[string]*sameDigestMapReply)
	(*output).requiredNumberOfMatches = requiredNumberOfMatches

	return output
}

func (obj *MapReplyRegistry) Add(replyDigest string, replyNodeId int, dataSize map[int]int) bool {

	reply := (*obj).registry[replyDigest]

	if reply == nil {

		reply = new(sameDigestMapReply)

		(*reply).sameDigestReplyNodeIds = make([]int, 0)
		(*reply).dataSize = dataSize

		(*obj).registry[replyDigest] = reply
	}

	(*reply).sameDigestReplyNodeIds = append((*reply).sameDigestReplyNodeIds, replyNodeId)

	if len((*reply).sameDigestReplyNodeIds) == (*obj).requiredNumberOfMatches {
		(*obj).mostMatchedWorkerReplyDigest = replyDigest
		return true
	}

	return false
}

func (obj *MapReplyRegistry) GetMostMatchedReply() (string, []int, map[int]int) {

	mostMatchedReply := (*obj).registry[(*obj).mostMatchedWorkerReplyDigest]

	return (*obj).mostMatchedWorkerReplyDigest, mostMatchedReply.sameDigestReplyNodeIds, mostMatchedReply.dataSize
}
