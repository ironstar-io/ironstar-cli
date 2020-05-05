package version

import (
	"fmt"
)

// Display - Pull the current version info and display in the console
func Display() {
	info := Get()

	fmt.Println(`
Ironstar CLI: ` + info.Version + `
Build Date:   ` + info.BuildDate + `
Compiler:     ` + info.GoVersion + `
Platform:     ` + info.Platform + `
	`)
}
