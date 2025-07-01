package funcs

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	source := []int{1, 6, 8, 2}
	res := Map(
		source, func(num int) string {
			return strconv.Itoa(num)
		},
	)
	assert.Equal(t, []string{"1", "6", "8", "2"}, res)
}

func TestFilter(t *testing.T) {
	source := []int{1, 6, 8, 2}
	res := Filter(
		source, func(num int, _ int) bool {
			return num > 5
		},
	)
	assert.Equal(t, []int{6, 8}, res)
}

func TestGroupBy(t *testing.T) {
	type item struct {
		kind string
		code int
	}
	source := []item{
		{
			kind: "A",
			code: 1,
		},
		{
			kind: "A",
			code: 2,
		},
		{
			kind: "B",
			code: 3,
		},
		{
			kind: "A",
			code: 4,
		},
		{
			kind: "C",
			code: 5,
		},
		{
			kind: "B",
			code: 6,
		},
	}
	res := GroupBy(
		source, func(item item) string {
			return item.kind
		},
	)
	assert.Equal(t, 3, len(res["A"]))
	assert.Equal(t, 2, len(res["B"]))
	assert.Equal(t, 1, len(res["C"]))
}

func TestAssociateBy(t *testing.T) {
	type item struct {
		kind string
		code int
	}
	source := []item{
		{
			kind: "A",
			code: 1,
		},
		{
			kind: "A",
			code: 2,
		},
		{
			kind: "B",
			code: 3,
		},
		{
			kind: "A",
			code: 4,
		},
		{
			kind: "C",
			code: 5,
		},
		{
			kind: "B",
			code: 6,
		},
	}
	res := AssociateBy(
		source, func(item item) string {
			return item.kind
		},
	)
	assert.Equal(t, 4, res["A"].code)
	assert.Equal(t, 6, res["B"].code)
	assert.Equal(t, 5, res["C"].code)
}

func TestChunk(t *testing.T) {
	source := []int{1, 6, 8, 2}
	res := Chunk(source, 3)
	assert.Equal(t, 2, len(res))
	assert.Equal(t, []int{1, 6, 8}, res[0])
	assert.Equal(t, []int{2}, res[1])
}

func TestReduce(t *testing.T) {
	source := []int{1, 6, 8, 2}
	res := Reduce(
		source, func(agg int, item int, _ int) int {
			return agg + item
		}, 0,
	)
	assert.Equal(t, 17, res)
}
