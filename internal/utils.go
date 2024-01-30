package internal

import (
	"fmt"
	"math"
)

func formatBytes(size float64) string {
	const unit = 1024

	if size < unit {
		return fmt.Sprintf("%.2f", size) + " B"
	}

	div := int(math.Log(size) / math.Log(unit))
	exp := "KMGTPE"[div-1]
	return fmt.Sprintf("%.2f ", size/math.Pow(unit, float64(div))) + string(exp) + "B"
}
