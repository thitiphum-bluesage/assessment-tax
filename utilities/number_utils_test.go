package utilities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatToTwoDecimals(t *testing.T) {
    tests := []struct {
        name     string
        input    float64
        expected float64
    }{
        {"Round down", 123.456, 123.46},
        {"Round up", 123.654, 123.65},
        {"Negative round down", -123.456, -123.46},
        {"Negative round up", -123.654, -123.65},
        {"No rounding needed", 123.45, 123.45},
        {"Zero", 0, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := FormatToTwoDecimals(tt.input)
            assert.Equal(t, tt.expected, result, "Output should match expected value")
        })
    }
}
