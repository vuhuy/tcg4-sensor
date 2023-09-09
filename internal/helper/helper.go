package helper

import (
	"fmt"
	"os"

	"github.com/vuhuy/tcg4-sensor/internal/global"
)

// Lousy method to convert hertz to milliseconds.
func ConvertHzToMilliseconds(frequency float64) uint {
	return uint(1.0 / frequency * 1000.0)
}

// Print version.
func PrintVersion(version string) {
	if *global.Version {
		fmt.Println(version)
		os.Exit(0)
	}
}
