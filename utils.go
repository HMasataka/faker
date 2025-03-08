package faker

import (
	"fmt"
	"strings"

	"github.com/HMasataka/perseus"
)

func BuildQuestionMarks(numRecords, numColumns int) string {
	r := fmt.Sprintf("(%v)", strings.Join(perseus.Repeat(numColumns, "?"), ","))

	if numRecords == 1 {
		return r
	}

	return strings.Join(perseus.Repeat(numRecords, r), ",")
}
