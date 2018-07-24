package hash

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

const separator = "."

// Calculate will return a hash according to the data given
func Calculate(data []string) (res string) {
	str := strings.Join(data, separator)
	sum := sha256.Sum256([]byte(str))
	res = fmt.Sprintf("%x", sum)
	return
}
