package mapper_test

import (
	"strconv"
	"testing"

	"github.com/4ND3R50N/go-tools/converter"
	"github.com/4ND3R50N/go-tools/mapper"

	"github.com/stretchr/testify/assert"
)

type A struct {
	ID *int
}

type FlatA struct {
	name   *string
	IDList []int
}

type B struct {
	ID *string
}

func TestMapper_Map(t *testing.T) {
	listA := []A{
		{
			ID: converter.ConvertToPointer(1),
		},
		{
			ID: converter.ConvertToPointer(2),
		},
	}
	listB := mapper.Map(listA, func(fromEntry A) B {
		return B{ID: converter.ConvertToPointer(strconv.Itoa(*fromEntry.ID))}
	})

	// test that listA did not get mutated
	assert.Equal(t, converter.ConvertToPointer(1), listA[0].ID)
	assert.Equal(t, converter.ConvertToPointer(2), listA[1].ID)
	assert.Equal(t, 2, len(listB))
	assert.Equal(t, converter.ConvertToPointer("1"), listB[0].ID)
	assert.Equal(t, converter.ConvertToPointer("2"), listB[1].ID)
}

func TestMapper_FlatMap(t *testing.T) {
	listA := []FlatA{
		{
			name: converter.ConvertToPointer("ID Collection 1"),
			IDList: []int{
				1,
				2,
			},
		},
		{
			name: converter.ConvertToPointer("ID Collection 2"),
			IDList: []int{
				100,
			},
		},
	}
	listB := mapper.FlatMap(listA, func(fromEntry FlatA) []B {
		var listOfB []B
		for _, id := range fromEntry.IDList {
			listOfB = append(listOfB, B{ID: converter.ConvertToPointer(strconv.Itoa(id))})
		}
		return listOfB
	})
	assert.Equal(t, converter.ConvertToPointer("ID Collection 1"), listA[0].name)
	assert.Equal(t, 1, listA[0].IDList[0])
	assert.Equal(t, 2, listA[0].IDList[1])
	assert.Equal(t, converter.ConvertToPointer("ID Collection 2"), listA[1].name)
	assert.Equal(t, 100, listA[1].IDList[0])
	assert.Equal(t, 3, len(listB))
	assert.Equal(t, converter.ConvertToPointer("1"), listB[0].ID)
	assert.Equal(t, converter.ConvertToPointer("2"), listB[1].ID)
	assert.Equal(t, converter.ConvertToPointer("100"), listB[2].ID)
}

func TestMapper_MapWithErr(t *testing.T) {
	type TestCase struct {
		name          string
		sliceToMap    []B
		targetSlice   []*A
		expectedError error
	}
	tests := []TestCase{
		{
			name: "mapping works",
			sliceToMap: []B{
				{
					ID: converter.ConvertToPointer("1"),
				},
				{
					ID: converter.ConvertToPointer("2"),
				},
			},
			targetSlice: []*A{
				{
					ID: converter.ConvertToPointer(1),
				},
				{
					ID: converter.ConvertToPointer(2),
				},
			},
			expectedError: nil,
		},
		{
			name: "mapping throws error",
			sliceToMap: []B{
				{
					ID: converter.ConvertToPointer("A"),
				},
				{
					ID: converter.ConvertToPointer("B"),
				},
			},
			expectedError: &strconv.NumError{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			listA, err := mapper.MapWithErr(tt.sliceToMap, func(fromEntry B) (*A, error) {
				value, err := strconv.Atoi(*fromEntry.ID)
				if err != nil {
					return nil, err
				}
				return &A{ID: &value}, nil
			})

			if tt.expectedError != nil {
				assert.ErrorAs(t, err, &tt.expectedError)
			}
			if tt.targetSlice != nil {
				assert.Equal(t, converter.ConvertToPointer("1"), tt.sliceToMap[0].ID)
				assert.Equal(t, converter.ConvertToPointer("2"), tt.sliceToMap[1].ID)
				assert.Equal(t, tt.targetSlice, listA)
			}
		})
	}
}
