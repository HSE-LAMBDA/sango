package src

import (
	"io/ioutil"
	"os"
)

func ParseFileAndUnmarshal(filename string) []byte {

	jsonFile, err := os.Open(filename)
	if err != nil {
		panic("File error")
	}

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic("Read error")
	}

	return bytes
}
