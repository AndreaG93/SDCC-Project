package replyregister

type Register struct {
	registry                      map[string]*sameDigestOutputProcessesReply
	requiredNumberOfMatches       int
	mostMatchedProcessReplyDigest string
}

type sameDigestOutputProcessesReply struct {
	processIDs     []int
	additionalData interface{}
}

func New(requiredNumberOfMatches int) *Register {

	output := new(Register)

	(*output).registry = make(map[string]*sameDigestOutputProcessesReply)
	(*output).requiredNumberOfMatches = requiredNumberOfMatches

	return output
}

func (obj *Register) AddReplyCheckingRequiredMatches(digest string, processId int, data interface{}) bool {

	reply := (*obj).registry[digest]

	if reply == nil {

		reply = new(sameDigestOutputProcessesReply)

		(*reply).processIDs = make([]int, 1)
		(*reply).processIDs[0] = processId
		(*reply).additionalData = data

		(*obj).registry[digest] = reply

	} else {
		(*reply).processIDs = append((*reply).processIDs, processId)
	}

	if len((*reply).processIDs) == (*obj).requiredNumberOfMatches {
		(*obj).mostMatchedProcessReplyDigest = digest
		return true
	}

	return false
}

func (obj *Register) GetMostMatchedReply() (string, []int, interface{}) {

	mostMatchedReply := (*obj).registry[(*obj).mostMatchedProcessReplyDigest]

	return (*obj).mostMatchedProcessReplyDigest, (*mostMatchedReply).processIDs, (*mostMatchedReply).additionalData
}
