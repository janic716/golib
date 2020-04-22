package file

import (
	"io/ioutil"
	"os"
	"path"
)

func ReadText(filename string) (contents string, err error) {
	data, err := ioutil.ReadFile(filename)
	if err == nil {
		contents = string(data)
	}
	return
}

func Remove(filename string) (err error) {
	return os.Remove(filename)
}

func WriteText(filename, content string) (int, error) {
	return WriteBinary(filename, []byte(content))
}

func WriteBinary(filename string, bytes []byte) (int, error) {
	os.MkdirAll(path.Dir(filename), os.ModePerm)
	file, err := os.Create(filename)
	if err != nil {
		return 0, err
	}
	defer func() {
		err = file.Close()
	}()
	return file.Write(bytes)
}
