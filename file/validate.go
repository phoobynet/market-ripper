package file

import (
	"errors"
	"fmt"
	"os"
)

func MustExist(filePath string) {
	if _, err := os.Stat(filePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			panic(fmt.Errorf("file \"%s\" does not exist", filePath))
		} else {
			panic(err)
		}
	}
}
