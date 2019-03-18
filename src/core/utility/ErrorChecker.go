package utility

const (
	InvalidInput = "ERROR: Invalid input"
)

func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}
