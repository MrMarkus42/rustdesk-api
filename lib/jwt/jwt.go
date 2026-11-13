package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Jwt struct {
	Key                 []byte
	TokenExpireDuration time.Duration
}

type UserClaims struct {
	UserId uint `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJwt(key string, tokenExpireDuration time.Duration) *Jwt {
	return &Jwt{
		Key:                 []byte(key),
		TokenExpireDuration: tokenExpireDuration,
	}
}

func (s *Jwt) GenerateToken(userId uint) string {
	return s.GenerateTokenWithExpire(userId, s.TokenExpireDuration)
}

// GenerateTokenWithExpire signs a token with its own lifetime (e.g. short-lived web client tokens).
func (s *Jwt) GenerateTokenWithExpire(userId uint, expire time.Duration) string {
	if len(s.Key) == 0 {
		fmt.Println("jwt key is nil")
		return ""
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		UserClaims{
			UserId: userId,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			},
		})
	token, err := t.SignedString(s.Key)
	if err != nil {
		fmt.Printf("jwt token generate error: %v", err)
		return ""
	}
	return token
}

func (s *Jwt) ParseToken(tokenString string) (uint, error) {
	return s.parse(tokenString)
}

// ParseTokenSignature checks the signature but not the expiry, for callers that keep
// the validity of the token elsewhere (the user_tokens table for client logins).
func (s *Jwt) ParseTokenSignature(tokenString string) (uint, error) {
	return s.parse(tokenString, jwt.WithoutClaimsValidation())
}

func (s *Jwt) parse(tokenString string, opts ...jwt.ParserOption) (uint, error) {
	opts = append(opts, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.Key, nil
	}, opts...)
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims.UserId, nil
	}
	return 0, err
}
