package filter_test

import (
	"strings"
	"testing"

	"github.com/4ND3R50N/go-tools/converter"
	"github.com/4ND3R50N/go-tools/filter"

	"github.com/stretchr/testify/assert"
)

func TestFilter_Filter(t *testing.T) {
	listOfElements := []*string{
		converter.ToPointer("Hello"), converter.ToPointer("World"),
	}
	matchingElements := filter.Filter(listOfElements, func(element *string) bool {
		return strings.Contains(*element, "r")
	})

	assert.Equal(t, converter.ToPointer("Hello"), listOfElements[0])
	assert.Equal(t, converter.ToPointer("World"), listOfElements[1])
	assert.Equal(t, 1, len(matchingElements))
	assert.Equal(t, converter.ToPointer("World"), matchingElements[0])
}

func TestFilter_Distinct(t *testing.T) {
	type TestCase struct {
		name           string
		input          []int
		expectedOutput []int
	}
	testCases := []TestCase{{
		name:           "remove duplicates",
		input:          []int{1, 1, 2, 1, 0, 2, 1, 2, 3},
		expectedOutput: []int{1, 2, 0, 3},
	}, {
		name:           "do nothing if there are no duplicates",
		input:          []int{1, 2, 0, 3},
		expectedOutput: []int{1, 2, 0, 3},
	}, {
		name:           "do nothing if slice is empty",
		input:          []int{},
		expectedOutput: []int{},
	}}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			output := filter.Distinct(tt.input)
			assert.Equal(t, output, tt.expectedOutput)
		})
	}
}
