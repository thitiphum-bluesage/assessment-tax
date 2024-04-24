package utilities

import "math"

func FormatToTwoDecimals(num float64) float64 {
    return math.Round(num*100)/100
}