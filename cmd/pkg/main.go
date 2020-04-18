package pkg

import (
	"github.com/spf13/cobra"
)

// PackageCmd - `iron package`
var PackageCmd = &cobra.Command{
	Use:   "package",
	Short: "",
	Long:  "",
	Run:   create,
}

// PkgCmd - `iron pkg`
var PkgCmd = &cobra.Command{
	Hidden: true,
	Use:    "pkg",
	Short:  "",
	Long:   "",
	Run:    create,
}
