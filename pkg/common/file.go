package common

import (
	"io/ioutil"
	"log"
	"os"
)

func CreateDir(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil && !os.IsExist(err) {
			return err
		}
	}

	return nil
}

func GetDirFileNum(path string) int {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil && !os.IsExist(err) {
			return 0
		}
	}
	fileInfoList, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
		return 0
	}

	return len(fileInfoList)
}
