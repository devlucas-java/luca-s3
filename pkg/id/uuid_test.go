package id_test

import (
	"testing"

	"github.com/devlucas-java/klyp-shop/pkg/id"
	"github.com/stretchr/testify/assert"
)

func TestNewUUID_IsUnique(t *testing.T) {
	a := id.NewUUID()
	b := id.NewUUID()
	assert.NotEqual(t, a, b)
}

func TestNewUUID_IsNotZero(t *testing.T) {
	u := id.NewUUID()
	var zero id.UUID
	assert.NotEqual(t, zero, u)
}

func TestParse_Valid(t *testing.T) {
	original := id.NewUUID()
	parsed, err := id.Parse(original.String())
	assert.NoError(t, err)
	assert.Equal(t, original, parsed)
}

func TestParse_Invalid(t *testing.T) {
	_, err := id.Parse("not-a-uuid")
	assert.Error(t, err)
}

func TestParse_Empty(t *testing.T) {
	_, err := id.Parse("")
	assert.Error(t, err)
}
