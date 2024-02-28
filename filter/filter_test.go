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
