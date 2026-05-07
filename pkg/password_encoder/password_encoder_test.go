package password_encoder_test

import (
	"testing"

	"github.com/devlucas-java/klyp-shop/pkg/password_encoder"
	"github.com/stretchr/testify/assert"
)

func TestEncoder_ReturnsHash(t *testing.T) {
	hash, err := password_encoder.Encoder("mypassword")
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, "mypassword", hash)
}

func TestEncoder_DifferentSaltsPerCall(t *testing.T) {
	h1, _ := password_encoder.Encoder("samepassword")
	h2, _ := password_encoder.Encoder("samepassword")
	assert.NotEqual(t, h1, h2)
}

func TestMatch_CorrectPassword(t *testing.T) {
	hash, err := password_encoder.Encoder("correctpass")
	assert.NoError(t, err)

	match, err := password_encoder.Match("correctpass", hash)
	assert.NoError(t, err)
	assert.True(t, match)
}

func TestMatch_WrongPassword(t *testing.T) {
	hash, err := password_encoder.Encoder("correctpass")
	assert.NoError(t, err)

	match, err := password_encoder.Match("wrongpass", hash)
	assert.NoError(t, err)
	assert.False(t, match)
}

func TestMatch_InvalidHash(t *testing.T) {
	_, err := password_encoder.Match("password", "not-a-valid-hash")
	assert.Error(t, err)
}

func TestMatch_EmptyPassword(t *testing.T) {
	hash, err := password_encoder.Encoder("somepass")
	assert.NoError(t, err)

	match, err := password_encoder.Match("", hash)
	assert.NoError(t, err)
	assert.False(t, match)
}
