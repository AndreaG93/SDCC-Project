package file_transmission

import "testing"

func Test_SFTP(t *testing.T) {

	StartSFTPServer()

	SendFile("localhost")

}
