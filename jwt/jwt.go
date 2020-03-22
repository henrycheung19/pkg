package jwt

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

// ErrTokenInvalid is returned when the token is invalid after parsed.
var ErrTokenInvalid = errors.New("ErrTokenInvalid")

// ExamTokenClaims contains the JWT claims of exam token.
type ExamTokenClaims struct {
	ExamPaperID   int `json:"exam_paper_id"`
	ExamSessionID int `json:"exam_session_id"`
	UserGroupID   int `json:"user_group_id"`
	jwt.StandardClaims
}

// IssueToken issues a new JWToken with the given claims and hmac secret.
func IssueToken(claims jwt.Claims, hmacSecret []byte) (string, error) {
	// Prepare JWToken
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Return signed token
	return token.SignedString(hmacSecret)
}

// Parse parses the jwt token into the given claim if it is valid.
func Parse(tokenStr string, claims *jwt.Claims, hmacSecret []byte) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenStr, *claims, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, ErrTokenInvalid
	}
	return token, nil
}
