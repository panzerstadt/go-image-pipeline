package shared

import (
	"os"
)

func Copy(from string, to string) error {
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
