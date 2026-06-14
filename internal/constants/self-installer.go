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

	// BinaryNameMacOS is the name of the Ironstar CLI macOS universal (amd64 + arm64) binary on Github
	BinaryNameMacOS = "iron-macos"

	// BinaryNameLinuxAMD64 is the name of the Ironstar CLI Linux x86-64 binary on Github
	BinaryNameLinuxAMD64 = "iron-linux-amd64"

	// BinaryNameLinuxARM64 is the name of the Ironstar CLI Linux arm64 binary on Github
	BinaryNameLinuxARM64 = "iron-linux-arm64"

	// BinaryNameWindows is the name of the Ironstar CLI Windows binary on Github
	BinaryNameWindows = "iron-windows.exe"

	// ActiveBinaryPathDarwin is the location of of the 'iron' command which is a symlink to the active Ironstar CLI version
	ActiveBinaryPathDarwin = "/usr/local/bin/iron"
)
