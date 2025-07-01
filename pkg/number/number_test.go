package number

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFirstThreeDigits(t *testing.T) {
	tests := []int{321310, 65, 0}
	res := make([]int, 0)
	for _, v := range tests {
		res = append(res, GetFirstThreeDigits(v))
	}
	assert.Equal(t, 321, res[0])
	assert.Equal(t, 65, res[1])
	assert.Equal(t, 0, res[2])
}
