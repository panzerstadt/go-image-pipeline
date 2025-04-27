package main

import "os"

func save(fullpath string, data []byte) error {
	err := os.WriteFile(fullpath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
