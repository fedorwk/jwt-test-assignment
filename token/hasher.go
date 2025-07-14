package token

import (
	"crypto/sha256"

	"golang.org/x/crypto/bcrypt"
)

type TokenHash []byte

type EncodedToken string

func (enc EncodedToken) String() string {
	return string(enc)
}

type Hasher interface {
	Hash(EncodedToken) (TokenHash, error)
}

func (enc EncodedToken) Bytes() []byte {
	return []byte(enc)
}

type BcryptHasher struct{}

func (h BcryptHasher) Hash(t EncodedToken) (TokenHash, error) {
	prehash := sha256.Sum256([]byte(t))
	hash, err := bcrypt.GenerateFromPassword(prehash[:], bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}
