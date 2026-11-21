package internal

import (
	"strconv"
	"strings"
)

// FormatVector преобразует срез float32 в строку формата pgvector: "[1.0,2.0,...]"
func FormatVector(v []float32) string {
	var sb strings.Builder
	sb.WriteByte('[')
	for i, val := range v {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(strconv.FormatFloat(float64(val), 'f', -1, 32))
	}
	sb.WriteByte(']')
	return sb.String()
}
