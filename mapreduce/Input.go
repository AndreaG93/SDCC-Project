package mapreduce

type Input interface {
	Split()
	getSplits() []Split
}
