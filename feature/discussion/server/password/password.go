package password

import (
	"math/rand/v2"
	"strings"
)

var charset = []rune("aAbBcCdDeEfFgGHhJjKkMmNnPpQqRrSsTtUuVvWwXxYyZz123456789~!@#$%^*")

func New() string {
	var sb strings.Builder

	for range 16 {
		sb.WriteRune(charset[rand.IntN(len(charset))])
	}

	return sb.String()
}
