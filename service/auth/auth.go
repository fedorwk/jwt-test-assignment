package auth

import (
	"context"
	"errors"
	"medods-auth/token"
	"medods-auth/user"
	"time"

	"github.com/google/uuid"
)

type AuthError error

var (
	ErrUserAgentChanged AuthError = errors.New("user-agent changed")
	ErrUserIDMissmatch  AuthError = errors.New("user id mismatch")

	ErrAccessTokenExpired  AuthError = errors.New("access token expired")
	ErrTokenExpired        AuthError = errors.New("token expired")
	ErrRefreshTokenExpired AuthError = errors.New("refresh token expired")

	ErrNilRefreshToken AuthError = errors.New("empty refresh token passed")

	ErrAccessTokenExpected  AuthError = errors.New("access token expected")
	ErrRefreshTokenExpected AuthError = errors.New("refresh token expected")

	ErrBlackListedToken AuthError = errors.New("blacklisted token provided")
)

type TokenPair struct {
	Access  *token.EncodedToken
	Refresh *token.EncodedToken
}

type RefreshTokenRecord struct {
	JTI       token.JTI
	User      user.User
	Hash      token.TokenHash
	CreatedAt time.Time

	RevokedAt *time.Time
}

type TokenHashRepository interface {
	Store(context.Context, *RefreshTokenRecord) error
	Get(context.Context, token.JTI) (*RefreshTokenRecord, error)
	DeleteByUserId(context.Context, uuid.UUID) error
}

type TokenBlackList interface {
	Add(context.Context, token.JTI) error
	Contains(context.Context, token.JTI) (bool, error)
}

type AuthService struct {
	refreshTokenRepo TokenHashRepository
	generator        token.Generator
	hasher           token.Hasher
	blacklist        TokenBlackList

	accessTTL  time.Duration
	refreshTTL time.Duration
	secret     []byte
}

type AuthServiceOptions struct {
	RefreshTokenRepo TokenHashRepository
	Generator        token.Generator
	Hasher           token.Hasher
	Blacklist        TokenBlackList

	Secret []byte

	AccessTTL  *time.Duration
	RefreshTTL *time.Duration
}

func NewAuthService(opts AuthServiceOptions) (*AuthService, error) {
	if opts.RefreshTokenRepo == nil {
		return nil, errors.New("nil token repository")
	}
	if opts.Generator == nil {
		return nil, errors.New("nil token generator")
	}
	if opts.Hasher == nil {
		return nil, errors.New("nil token hasher")
	}
	if opts.Blacklist == nil {
		return nil, errors.New("nil token blacklist repository")
	}
	if opts.Secret == nil {
		return nil, errors.New("nil signing secrect")
	}
	if opts.AccessTTL == nil {
		return nil, errors.New("nil access token ttl")
	}
	if opts.RefreshTTL == nil {
		return nil, errors.New("nil refresh token ttl")
	}
	return &AuthService{
		refreshTokenRepo: opts.RefreshTokenRepo,
		generator:        opts.Generator,
		hasher:           opts.Hasher,
		blacklist:        opts.Blacklist,

		accessTTL:  *opts.AccessTTL,
		refreshTTL: *opts.RefreshTTL,
		secret:     opts.Secret,
	}, nil
}

func (s *AuthService) GenerateTokens(u user.User) (TokenPair, error) {
	access := s.generator.Generate(token.Options{
		User: u,
		TTL:  s.accessTTL,
	})
	accessEnc, err := s.encodeToken(access)
	if err != nil {
		return TokenPair{}, err
	}

	refresh := s.generator.Generate(token.Options{
		User: u,
		TTL:  s.refreshTTL,
	})
	refreshEnc, err := s.encodeToken(refresh)
	if err != nil {
		return TokenPair{}, err
	}

	jti, err := refresh.JTI()
	if err != nil {
		return TokenPair{}, err
	}
	hash, err := s.hasher.Hash(refreshEnc)
	if err != nil {
		return TokenPair{}, err
	}
	tokenRecord := RefreshTokenRecord{
		JTI:       jti,
		User:      u,
		Hash:      hash,
		CreatedAt: time.Now(),
	}

	err = s.refreshTokenRepo.Store(context.TODO(), &tokenRecord)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		Access:  &accessEnc,
		Refresh: &refreshEnc,
	}, nil
}

func (s *AuthService) Refresh(u user.User, pair TokenPair) (TokenPair, error) {
	refresh, err := s.decodeToken(*pair.Refresh)
	if err != nil {
		return TokenPair{}, nil
	}
	err = s.Validate(&u, refresh)
	if err != nil {
		return TokenPair{}, err
	}

	err = s.RevokeTokens(u, *pair.Access)
	if err != nil {
		return TokenPair{}, err
	}

	return s.GenerateTokens(u)
}

func (s *AuthService) RevokeTokens(u user.User, access token.EncodedToken) error {
	decoded, err := s.decodeToken(access)
	if err != nil {
		return err
	}
	err = s.Validate(&u, decoded)
	if err != nil {
		return err
	}
	// revoke access
	err = s.revokeAccessToken(decoded)
	if err != nil {
		return err
	}

	err = s.refreshTokenRepo.DeleteByUserId(context.TODO(), u.Id)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) Validate(u *user.User, t *token.Token) error {
	if exp, err := t.Expires(); err != nil {
		return err
	} else if time.Now().After(exp) {
		return ErrTokenExpired
	}

	claims, err := t.GetClaims()
	if err != nil {
		return err
	}
	if u != nil {
		if claims.UserAgent != u.UserAgent {
			return ErrUserAgentChanged
		}

		if userid, err := t.UserID(); err != nil {
			return err
		} else if userid.String() != u.Id.String() {
			return ErrUserIDMissmatch
		}
	}

	jti, err := t.JTI()
	if err != nil {
		return err
	}
	blacklisted, err := s.blacklist.Contains(context.TODO(), jti)
	if err != nil {
		return err
	}
	if blacklisted {
		return ErrBlackListedToken
	}

	return nil
}

func (s *AuthService) ExtractUserID(enc *token.EncodedToken) (uuid.UUID, error) {
	token, err := s.generator.Decode(enc.String(), s.secret)
	if err != nil {
		return uuid.Nil, err
	}
	err = s.Validate(nil, token)
	if err != nil {
		return uuid.Nil, err
	}
	id, err := token.UserID()
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (s *AuthService) encodeToken(t *token.Token) (token.EncodedToken, error) {
	enc, err := s.generator.Encode(t, s.secret)
	if err != nil {
		return "", err
	}
	return token.EncodedToken(enc), nil
}

func (s *AuthService) decodeToken(enc token.EncodedToken) (*token.Token, error) {
	token, err := s.generator.Decode(enc.String(), s.secret)
	if err != nil {
		return nil, err
	}
	return token, err
}

func (s *AuthService) revokeAccessToken(t *token.Token) error {
	jti, err := t.JTI()
	if err != nil {
		return err
	}
	return s.blacklist.Add(context.TODO(), jti)
}
