package main

import (
	"fmt"
	"os"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func main() {
	pathStr := "./images"
	res, err := PathExists(pathStr)
	if err != nil {
		fmt.Println("err", err)
	}
	if res == false {
		err = os.Mkdir(pathStr, 0777)
		if err != nil {
			fmt.Println("mkdir error", err)
		}
	}
}
