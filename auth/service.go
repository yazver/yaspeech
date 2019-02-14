package auth

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// ServiceAccount ...
type ServiceAccount struct {
	ID         string
	KeyID      string
	PrivateKey string
}

// JWT generation.
func (a *ServiceAccount) signedToken() (string, error) {
	issuedAt := time.Now()
	token := jwt.NewWithClaims(ps256WithSaltLengthEqualsHash, jwt.StandardClaims{
		Issuer:    a.ID,
		IssuedAt:  issuedAt.Unix(),
		ExpiresAt: issuedAt.Add(time.Hour).Unix(),
		Audience:  "https://iam.api.cloud.yandex.net/iam/v1/tokens",
	})
	token.Header["kid"] = a.KeyID

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(a.PrivateKey))
	if err != nil {
		return "", err
	}
	signed, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return signed, nil
}

// By default, the Go RSA-PSS algorithm uses PSSSaltLengthAuto,
// but https://tools.ietf.org/html/rfc7518#section-3.5 says that
// the size of the salt value should be the same size as the hash function output.
// After fixing https://github.com/dgrijalva/jwt-go/issues/285,
// it can be replaced with jwt.SigningMethodPS256
var ps256WithSaltLengthEqualsHash = &jwt.SigningMethodRSAPSS{
	SigningMethodRSA: jwt.SigningMethodPS256.SigningMethodRSA,
	Options: &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthEqualsHash,
	},
}

// GetToken return IAM token
func (a *ServiceAccount) GetToken() (string, error) {
	jot, err := a.signedToken()
	if err != nil {
		return "", err
	}
	resp, err := http.Post(
		"https://iam.api.cloud.yandex.net/iam/v1/tokens",
		"application/json",
		strings.NewReader(fmt.Sprintf(`{"jwt":"%s"}`, jot)),
	)
	if err != nil {
		return "", err
	}

	return processIAMResponse(resp)
}
