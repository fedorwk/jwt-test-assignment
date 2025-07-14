package token

import (
	"errors"
	"medods-auth/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	ClaimExpireAt = "exp"
	ClaimIssuedAt = "iat"
	ClaimSubject  = "sub"
	ClaimJWTID    = "jti"

	ClaimUserAgent = "user_agent"
)

var ErrUnexpextedClaimType = errors.New("unexpected type for claims")
var ErrNoClaimsInToken = errors.New("no claims in decoded token")
var ErrParsingTokenId = errors.New("err parsing token id")

type tokenType string

const (
	TokenTypeUnknown tokenType = "unknown"
	TokenTypeAccess  tokenType = "access"
	TokenTypeRefresh tokenType = "refresh"
)

type Token struct {
	t      *jwt.Token
	claims *Claims
}

type Claims struct {
	jwt.RegisteredClaims
	UserAgent string
	TokenType tokenType
}

type JTI = uuid.UUID

func (t *Token) Expires() (time.Time, error) {
	exp, err := t.t.Claims.GetExpirationTime()
	if err != nil {
		return time.Time{}, err
	}
	return exp.Time, nil
}

func (t *Token) JTI() (JTI, error) {
	if t.claims == nil {
		return JTI(uuid.Nil), ErrNoClaimsInToken
	}
	id, err := uuid.Parse(t.claims.ID)
	if err != nil {
		return JTI(uuid.Nil), ErrParsingTokenId
	}
	return JTI(id), nil
}

func (t *Token) UserID() (uuid.UUID, error) {
	if t.claims == nil {
		return uuid.Nil, ErrNoClaimsInToken
	}
	id, err := uuid.Parse(t.claims.Subject)
	if err != nil {
		return uuid.Nil, ErrParsingTokenId
	}
	return id, nil
}

func (t *Token) Type() (tokenType, error) {
	if t.claims == nil {
		return TokenTypeUnknown, ErrNoClaimsInToken
	}

	return t.claims.TokenType, nil
}

func (t *Token) GetClaims() (*Claims, error) {
	return t.extractClaims()
}

func (t *Token) extractClaims() (*Claims, error) {
	if claims, ok := t.t.Claims.(*Claims); ok {
		return claims, nil
	}
	return nil, ErrUnexpextedClaimType
}

type Options struct {
	User user.User
	TTL  time.Duration
	Type tokenType
}

type Generator interface {
	Generate(Options) *Token
	Encode(token *Token, secret []byte) (string, error)
	Decode(token string, secret []byte) (*Token, error)
}

type SHA512Generator struct{}

func (g *SHA512Generator) Generate(opts Options) *Token {
	now := time.Now()

	c := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   opts.User.Id.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(opts.TTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        generateJTI().String(),
		},
		UserAgent: opts.User.UserAgent,
		TokenType: opts.Type,
	}

	return &Token{
		t:      jwt.NewWithClaims(jwt.SigningMethodHS512, c),
		claims: c,
	}
}

func (g *SHA512Generator) Encode(t *Token, secret []byte) (string, error) {
	return t.t.SignedString(secret)
}

func (g *SHA512Generator) Decode(t string, secret []byte) (*Token, error) {
	decoded, err := jwt.ParseWithClaims(
		t,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			return secret, nil
		},
		jwt.WithValidMethods([]string{
			jwt.SigningMethodHS512.Alg(),
		}),
	)
	if err != nil {
		return nil, err
	}

	token := &Token{
		t: decoded,
	}
	token.claims, err = token.extractClaims()
	if err != nil {
		return nil, err
	}
	return token, nil
}

// FEAT: DI, Typing
func generateJTI() JTI {
	return uuid.New()
}
