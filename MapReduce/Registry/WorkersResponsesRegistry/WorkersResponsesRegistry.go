package WorkersResponsesRegistry

type WorkersResponsesRegistry struct {
	workersResponses                 map[string][]string
	requiredNumberOfMatches          int
	isRequiredNumberOfMatchesReached chan bool
	mostMatchedWorkerResponseDigest  string
}

func New(requiredNumberOfMatches int, isRequiredNumberOfMatchesReached chan bool) *WorkersResponsesRegistry {

	output := new(WorkersResponsesRegistry)

	(*output).workersResponses = make(map[string][]string)
	(*output).requiredNumberOfMatches = requiredNumberOfMatches
	(*output).isRequiredNumberOfMatchesReached = isRequiredNumberOfMatchesReached

	return output
}

func (obj *WorkersResponsesRegistry) AddWorkerResponse(digest string, workerAddress string) bool {

	workerAddresses := (*obj).workersResponses[digest]

	if workerAddresses == nil {
		workerAddresses = make([]string, 0)
	}

	(*obj).workersResponses[digest] = append(workerAddresses, workerAddress)

	if len((*obj).workersResponses[digest]) == (*obj).requiredNumberOfMatches {
		(*obj).mostMatchedWorkerResponseDigest = digest
		(*obj).isRequiredNumberOfMatchesReached <- true
		return true
	}

	return false
}

func (obj *WorkersResponsesRegistry) GetMostMatchedWorkerResponse() (string, []string) {
	return (*obj).mostMatchedWorkerResponseDigest, (*obj).workersResponses[(*obj).mostMatchedWorkerResponseDigest]
}
