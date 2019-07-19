package Data

type Split interface {
	PerformTask() (string, []byte, error)
}
