package mapreduce

type Split interface {
	performMapTask()
}
