package console

import (
	"fmt"
	"os"
	"runtime"

	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"github.com/spf13/viper"
)

// Println ...
func Println(message string, replacement string) {
	if viper.Get("Enableemoji") == false || runtime.GOOS == "windows" {
		fmt.Println(StripEmoji(message, replacement))
	} else {
		fmt.Println(message + "  ")
	}
}

// Printf ...
func Printf(message string, replacement string) {
	if viper.Get("Enableemoji") == false || runtime.GOOS == "windows" {
		fmt.Printf(StripEmoji(message, replacement))
	} else {
		fmt.Printf(message + "  ")
	}
}

// SpinStart ...
func SpinStart(message string) *wow.Wow {
	if viper.Get("Enableemoji") == false || runtime.GOOS == "windows" {
		Println(message, "")
		return nil
	}

	wo := wow.New(os.Stdout, spin.Get(spin.Dots), `   `+message)
	wo.Start()

	return wo
}

// SpinPersist ...
func SpinPersist(wo *wow.Wow, emoji string, message string) {
	if viper.Get("Enableemoji") == false || runtime.GOOS == "windows" {
		Println(message, "")
	} else {
		wo.PersistWith(spin.Spinner{Frames: []string{emoji}}, `  `+message)
	}
}
