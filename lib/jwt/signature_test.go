package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestParseTokenSignatureIgnoresExpiry(t *testing.T) {
	s := NewJwt("test-key", time.Hour)
	expired := s.GenerateTokenWithExpire(42, -time.Minute)

	if uid, _ := s.ParseToken(expired); uid != 0 {
		t.Errorf("ParseToken accepted an expired token (uid %d)", uid)
	}
	if uid, err := s.ParseTokenSignature(expired); uid != 42 {
		t.Errorf("ParseTokenSignature = %d, %v; want 42", uid, err)
	}
}

func TestParseTokenSignatureRejectsForeignTokens(t *testing.T) {
	s := NewJwt("test-key", time.Hour)
	other := NewJwt("other-key", time.Hour).GenerateToken(42)
	if uid, _ := s.ParseTokenSignature(other); uid != 0 {
		t.Errorf("accepted a token signed with another key (uid %d)", uid)
	}

	none, err := jwt.NewWithClaims(jwt.SigningMethodNone, UserClaims{UserId: 42}).
		SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatal(err)
	}
	if uid, _ := s.ParseTokenSignature(none); uid != 0 {
		t.Errorf("accepted an unsigned token (uid %d)", uid)
	}

	hs512, err := jwt.NewWithClaims(jwt.SigningMethodHS512, UserClaims{UserId: 42}).
		SignedString([]byte("test-key"))
	if err != nil {
		t.Fatal(err)
	}
	if uid, _ := s.ParseTokenSignature(hs512); uid != 0 {
		t.Errorf("accepted a token with another algorithm (uid %d)", uid)
	}
}
