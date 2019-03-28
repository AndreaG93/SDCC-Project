package utility

const (
	InvalidInput = "ERROR: Invalid input"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
