package constants

const (
	// BaseBinaryURL is where Ironstar releases are available
	BaseBinaryURL = "https://github.com/ironstar-io/ironstar-cli/releases/download/"

	// BaseInstallPathLinux is where Ironstar CLI binaries are installed on Linux
	BaseInstallPathLinux = ".ironstar/bin"

	// BaseInstallPathDarwin is where Ironstar CLI binaries are installed on macOS/Darwin
	BaseInstallPathDarwin = ".ironstar/bin"

	// BaseInstallPathWindows is where Ironstar CLI binaries are installed on Windows
	BaseInstallPathWindows = "AppData/Local/Ironstar/CLI"

	// BinaryNameLinux is the name of the Ironstar CLI Linux binary on Github
	BinaryNameLinux = "iron-linux-amd64"

	// BinaryNameIntelMacOS is the name of the Ironstar CLI (Intel) macOS binary on Github
	BinaryNameIntelMacOS = "iron-macos"

	// BinaryNameARMMacOS is the name of the Ironstar CLI (ARM) macOS binary on Github
	BinaryNameARMMacOS = "iron-macos-arm64"

	// BinaryNameWindows is the name of the Ironstar CLI Windows binary on Github
	BinaryNameWindows = "iron-windows.exe"

	// ActiveBinaryPathDarwin is the location of of the 'iron' command which is a symlink to the active Ironstar CLI version
	ActiveBinaryPathDarwin = "/usr/local/bin/iron"
)
