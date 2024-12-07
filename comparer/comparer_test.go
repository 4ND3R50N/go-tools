package comparer_test

import (
	"testing"

	"github.com/4ND3R50N/go-tools/comparer"
	"github.com/4ND3R50N/go-tools/converter"
	"github.com/stretchr/testify/assert"
)

func TestComparer_Any(t *testing.T) {
	type TestCase struct {
		name           string
		slice          []*string
		sliceCopy      []*string
		expectedResult bool
	}
	tests := []TestCase{
		{
			name:           "return false for empty slice",
			expectedResult: false,
		}, {
			name:           "return false for slice not fulfilling predicate",
			slice:          []*string{converter.ToPointer("b")},
			sliceCopy:      []*string{converter.ToPointer("b")},
			expectedResult: false,
		}, {
			name:           "return true for slice fulfilling predicate",
			slice:          []*string{converter.ToPointer("b"), converter.ToPointer("a")},
			sliceCopy:      []*string{converter.ToPointer("b"), converter.ToPointer("a")},
			expectedResult: true,
		},
	}

	// test if slice contains "a" at least once
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			contains := comparer.Any(tt.slice, func(elem *string) bool {
				return *elem == "a"
			})
			assert.Equal(t, contains, tt.expectedResult)
			// assert that the original slice was not modified
			assert.Equal(t, tt.slice, tt.sliceCopy)
		})
	}
}

func TestComparer_All(t *testing.T) {
	type TestCase struct {
		name           string
		slice          []*string
		sliceCopy      []*string
		expectedResult bool
	}
	tests := []TestCase{
		{
			name:           "return true for empty slice",
			expectedResult: true,
		}, {
			name:           "return false for slice not fulfilling predicate",
			slice:          []*string{converter.ToPointer("a"), converter.ToPointer("b")},
			sliceCopy:      []*string{converter.ToPointer("a"), converter.ToPointer("b")},
			expectedResult: false,
		}, {
			name:           "return true for slice fulfilling predicate",
			slice:          []*string{converter.ToPointer("a"), converter.ToPointer("a")},
			sliceCopy:      []*string{converter.ToPointer("a"), converter.ToPointer("a")},
			expectedResult: true,
		},
	}

	// test if slice only contains "a"
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			contains := comparer.All(tt.slice, func(elem *string) bool {
				return *elem == "a"
			})
			assert.Equal(t, contains, tt.expectedResult)
			// assert that the original slice was not modified
			assert.Equal(t, tt.slice, tt.sliceCopy)
		})
	}
}
