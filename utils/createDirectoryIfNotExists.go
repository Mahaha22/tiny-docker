package utils

import (
	"fmt"
	"os"
)

func CreateDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			return err
		}
		fmt.Printf("Created directory: %s\n", path)
	} else {
		//fmt.Printf("Directory already exists: %s\n", path)
		return err
	}
	return nil
}
