package fs

import (
	"fmt"
	"log"
	"os/user"
)

// HomeDir - Return the users' $HOME dir
func HomeDir() string {
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("Ironstar CLI encountered a fatal error and had to stop")
		log.Fatal(err)
	}

	return usr.HomeDir
}
