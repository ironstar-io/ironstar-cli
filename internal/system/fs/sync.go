package fs

import (
	"os"
	"time"
)

// WaitForSync polls the specified file path locally until it appears
func WaitForSync(file string, retries int) (err error) {
	for i := 1; i <= retries; i++ {
		_, err = os.Stat(file)
		if err == nil {
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	return err
}
