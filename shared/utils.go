package shared

import (
	"math/rand/v2"
	"os"
	"time"
)

func RandDuration(min, max int) time.Duration {
	randNumber := rand.IntN(max-min+1) + min
	return time.Duration(randNumber) * time.Second
}

func FileCopy(from string, to string) error {
	contents, err := os.ReadFile(from)
	if err != nil {
		return err
	}

	err = os.WriteFile(to, contents, 0644)
	if err != nil {
		return err
	}
	return nil
}

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
