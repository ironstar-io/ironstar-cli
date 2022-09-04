package security

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/ironstar-io/ironstar-cli/internal/system/fs"
)

// CheckFilePermissions - Checks ~/.ironstar and ~/.ironstar/credentials.yml to ensure they have the correct permissions set
// See: https://ironstar.atlassian.net/browse/NANKAI-246
func CheckFilePermissions() {
	WarnPathPermissions(filepath.Join(fs.HomeDir(), ".ironstar"), "drwx------", "700")
	WarnPathPermissions(filepath.Join(fs.HomeDir(), ".ironstar", "credentials.yml"), "-r--------", "400")
}

func WarnPathPermissions(path, permission, octal string) {
	exists := fs.CheckExists(path)
	if exists {
		stat, err := os.Stat(path)
		if err != nil {
			return
		}

		if stat.Mode().String() != permission {
			fmt.Println()
			color.Yellow("WARNING!")
			color.Yellow(`The permissions for the file '` + path + `' are set to '` + stat.Mode().String() + `'. For security, Ironstar would highly recommend changing the permissions for this directory to '` + permission + `' with the commmand 'chmod ` + octal + ` ` + path + `'`)
		}
	}
}
