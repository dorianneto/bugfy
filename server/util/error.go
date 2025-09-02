package util

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

func GenerateFingerprint(message string) string {
	var parts []string
	parts = append(parts, message)

	data := strings.Join(parts, "|")
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)[:16]
}
