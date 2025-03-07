package faker

import (
	"strings"

	"github.com/HMasataka/perseus"
)

func BuildQuestionMarks(n int) string {
	return strings.Join(perseus.Repeat(n, "?"), ",")
}
