package ba

import "io/ioutil"

type FileReaderFn func(path string) (string, error)

func DefaultFileReader(path string) (string, error) {
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(buff), nil
}
