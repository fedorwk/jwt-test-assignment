package auth_test

import (
	"medods-auth/service/auth"
	"medods-auth/test/testutil"
	"medods-auth/token"
	"medods-auth/user"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var TestUser = user.User{
	Id:        uuid.New(),
	UserAgent: "User-Agent Test",
}

var TestUserAgentChanged = user.User{
	Id:        TestUser.Id,
	UserAgent: "User-Agent Changed",
}

func TestAuthService(t *testing.T) {
	assert := assert.New(t)

	testRepo := testutil.NewTestInmemoryRepo()
	defer testRepo.Close()
	defer testRepo.DumpContents()

	accessTTL := time.Second * 1
	refreshTTL := time.Second * 2

	service, err := auth.NewAuthService(auth.AuthServiceOptions{
		RefreshTokenRepo: testRepo,
		Blacklist:        testRepo,

		Generator: &token.SHA512Generator{},
		Hasher:    &token.BcryptHasher{},

		Secret:     []byte("test_secret"),
		AccessTTL:  &accessTTL,
		RefreshTTL: &refreshTTL,
	})
	assert.Nil(err)
	assert.NotNil(service)

	tokenPair, err := service.GenerateTokens(TestUser)
	assert.Nil(err)
	assert.NotNil(tokenPair.Access)
	assert.NotNil(tokenPair.Refresh)
	assert.NotEqual(tokenPair.Access, tokenPair.Refresh)

	_, err = service.Refresh(TestUserAgentChanged, tokenPair)
	assert.Equal(err, auth.ErrUserAgentChanged)

	updTokenPair, err := service.Refresh(TestUser, tokenPair)
	assert.Nil(err)
	assert.NotEqual(updTokenPair, tokenPair)

	_, err = service.Refresh(TestUser, tokenPair)
	assert.Equal(auth.ErrBlackListedToken, err, "old token should be revoken")
}
