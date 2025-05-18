package main

import "os"

func FileSave(fullpath string, data []byte) error {
	err := os.WriteFile(fullpath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func FileAppend(fullpath string, data []byte) error {
	f, err := os.OpenFile(fullpath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}
