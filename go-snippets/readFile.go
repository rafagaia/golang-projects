package snippets

import (
	"io/ioutil"
	"os"
)

func readFile(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	// ensures that the file will be closed when the function exits, regardless of whether an error occurs during reading.
	defer file.Close()

	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return contents, nil
}
