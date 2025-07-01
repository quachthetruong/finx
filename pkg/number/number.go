package number

import (
	"math"
)

func GetFirstThreeDigits(num int) int {
	// Calculate the number of digits in the input number
	numDigits := int(math.Log10(float64(num))) + 1
	if numDigits <= 3 {
		return num
	}
	divisor := int(math.Pow10(numDigits - 3))
	firstThreeDigits := num / divisor

	return firstThreeDigits
}
