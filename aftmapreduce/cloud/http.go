package cloud

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func Upload(url string, data []byte) error {

	var err error

	if request, err := http.NewRequest("PUT", url, bytes.NewReader(data)); err == nil {
		_, err = http.DefaultClient.Do(request)
	}

	return err
}

func Download(url string) ([]byte, error) {

	var err error

	if request, err := http.NewRequest("GET", url, nil); err == nil {
		if response, err := http.DefaultClient.Do(request); err == nil && response != nil {
			return ioutil.ReadAll(response.Body)
		}
	}

	return nil, err
}
