package util

import (
	"fmt"
	"os"
)

func ByteCountToReadable(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func CheckFileIsExist(name string) bool {
	exist := true
	if _, err := os.Stat(name); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
