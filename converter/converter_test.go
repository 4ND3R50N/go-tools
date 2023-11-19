package converter_test

import (
	"testing"

	"github.com/4ND3R50N/go-tools/converter"
	"github.com/stretchr/testify/assert"
)

func TestConverter_ConvertToPointer(t *testing.T) {
	normalString := "hey"
	pointerString := converter.ConvertToPointer(normalString)
	assert.Equal(t, &normalString, pointerString)
}
