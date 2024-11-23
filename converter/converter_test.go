package converter_test

import (
	"encoding/json"
	"github.com/4ND3R50N/go-tools/converter"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConverter_ToPointer(t *testing.T) {
	normalString := "hey"
	pointerString := converter.ToPointer(normalString)
	assert.Equal(t, &normalString, pointerString)
}

func TestConverter_ToValue(t *testing.T) {
	pointerString := converter.ToPointer("hey")
	normalString := converter.ToValue(pointerString)
	assert.Equal(t, normalString, *pointerString)
}

func TestConverter_ToPointerOrNil(t *testing.T) {
	normalString := "hey"
	pointerString := converter.ToPointerOrNil(normalString)
	assert.Equal(t, &normalString, pointerString)

	emptyString := ""
	assert.Nil(t, converter.ToPointerOrNil(emptyString))

	zeroInt := 0
	assert.Nil(t, converter.ToPointerOrNil(zeroInt))
}

func TestConverter_ToValueOrZero(t *testing.T) {
	pointerString := converter.ToPointer("hey")
	normalString := converter.ToValueOrZero(pointerString)
	assert.Equal(t, normalString, *pointerString)

	pointerString = nil
	assert.Equal(t, converter.ToValueOrZero(pointerString), "")

	var pointerInt *uint64
	assert.Equal(t, converter.ToValueOrZero(pointerInt), uint64(0))
}

func TestConverter_FromInterfaceTo(t *testing.T) {
	type TestType map[string]string
	// The real datatype.
	anyTestType := TestType{
		"test": "test",
	}
	// Convert anyTestType to interface{} but as string representation.
	testTypeBytes, err := json.Marshal(anyTestType)
	assert.NoError(t, err)
	var myInterface interface{} = string(testTypeBytes)

	// Convert '{"test":"test"}' interface data back to TestType.
	againTestType, err := converter.FromInterfaceTo[TestType](myInterface)
	assert.NoError(t, err)
	assert.Equal(t, anyTestType, againTestType)
}
