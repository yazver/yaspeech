package auth

import (
	"errors"
	"net/http"
	"strings"
)

// YandexAccount ...
type YandexAccount struct {
	OAuth string
}

// GetToken return IAM token
func (a *YandexAccount) GetToken() (string, error) {
	const url = "https://iam.api.cloud.yandex.net/iam/v1/tokens"

	resp, err := http.Post(url, "application/json", strings.NewReader(`{"yandexPassportOauthToken": "`+a.OAuth+`"}`))
	if err != nil {
		return "", errors.New("Unable to get token: " + err.Error())
	}

	return processIAMResponse(resp)
}
