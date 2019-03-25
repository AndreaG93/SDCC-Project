package utility

import "os"

func GenerateDirectory(path string) error {

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {

			return os.MkdirAll(path, os.ModePerm)

		} else {

			return nil
		}
	} else {
		return err
	}
}
