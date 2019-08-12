package reply

type ReduceReplyRegistry struct {
	registry                     map[string][]int
	requiredNumberOfMatches      int
	mostMatchedWorkerReplyDigest string
}

func NewReduceReplyRegistry(requiredNumberOfMatches int) *ReduceReplyRegistry {

	output := new(ReduceReplyRegistry)

	(*output).registry = make(map[string][]int)
	(*output).requiredNumberOfMatches = requiredNumberOfMatches

	return output
}

func (obj *ReduceReplyRegistry) Add(replyDigest string, replyNodeId int) bool {

	reply := (*obj).registry[replyDigest]

	if reply == nil {

		reply = make([]int, 0)

		(*obj).registry[replyDigest] = reply
	}

	reply = append(reply, replyNodeId)

	if len(reply) == (*obj).requiredNumberOfMatches {
		(*obj).mostMatchedWorkerReplyDigest = replyDigest
		return true
	}

	return false
}

func (obj *ReduceReplyRegistry) GetMostMatchedReply() (string, []int) {
	return (*obj).mostMatchedWorkerReplyDigest, (*obj).registry[(*obj).mostMatchedWorkerReplyDigest]
}
