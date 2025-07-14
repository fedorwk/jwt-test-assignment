package token

import (
	"medods-auth/user"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var testSecret = []byte("test_secret")

func TestGenerate(t *testing.T) {
	assert := assert.New(t)
	generator := &SHA512Generator{}

	testId := uuid.New()

	token := generator.Generate(Options{
		User: user.User{
			Id:        testId,
			UserAgent: "test_user",
		},
		TTL:  time.Second,
		Type: TokenTypeAccess,
	})

	encoding, err := generator.Encode(token, testSecret)
	assert.Nil(err, "err encoding token")

	parsed, err := generator.Decode(encoding, testSecret)
	assert.Nil(err, "err decoding token")
	assert.Equal(encoding, parsed.t.Raw)

	subject, err := parsed.t.Claims.GetSubject()
	assert.Nil(err, "getting subject claim")
	assert.Equal(subject, testId.String())

	assert.NotNil(parsed.claims)
	assert.Equal("test_user", parsed.claims.UserAgent)
	id, err := parsed.claims.GetSubject()
	assert.Nil(err)
	assert.Equal(testId.String(), id)
	assert.Equal(TokenTypeAccess, parsed.claims.TokenType)

	ttype, err := parsed.Type()
	assert.Nil(err)
	assert.Equal(TokenTypeAccess, ttype)
}
