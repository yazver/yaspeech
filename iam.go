package yaspeech

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

var token *IAMToken

// IAMToken is IAM token used for requests
type IAMToken struct {
	oauth      string
	err        error
	token      string
	lastUpdate time.Time
	mutex      sync.RWMutex
	ctx        context.Context
}

// InitIAMToken initialize IAM token.
func InitIAMToken(oauth string) {
	token = newIAMToken(nil, oauth)
}

// InitIAMTokenContext initialize IAM token
// The token will be updated in the background.
func InitIAMTokenContext(ctx context.Context, oauth string) {
	token = newIAMToken(ctx, oauth)
}

func newIAMToken(ctx context.Context, oauth string) *IAMToken {
	t := &IAMToken{ctx: ctx, oauth: oauth}
	t.init()
	return t
}

func (t *IAMToken) init() {
	if t.ctx != nil {
		go func() {
			t.update(true)
			ticker := time.NewTicker(time.Minute)
			for {
				select {
				case <-t.ctx.Done():
					return
				case <-ticker.C:
					t.update(true)
				}
			}
		}()
	}
}

// Get IAM token
func (t *IAMToken) Get() (string, error) {
	token.mutex.Lock()
	defer token.mutex.Unlock()

	if t.token == "" {
		t.update(false)
	}
	return t.token, t.err
}

func (t *IAMToken) update(lock bool) {
	if lock {
		t.mutex.Lock()
		defer t.mutex.Unlock()
	}

	if t.token != "" && time.Since(t.lastUpdate) < time.Hour*8 {
		return
	}

	token, err := t.requestToken()
	if err != nil && t.token != "" && time.Since(t.lastUpdate) < (time.Hour*11+time.Minute*50) {
		return
	}
	t.token, t.err, t.lastUpdate = token, err, time.Now()
}

func (t *IAMToken) requestToken() (string, error) {
	const url = "https://iam.api.cloud.yandex.net/iam/v1/tokens"

	resp, err := http.Post(url, "application/json", strings.NewReader(`{"yandexPassportOauthToken": "`+t.oauth+`"}`))
	if err != nil {
		return "", errors.New("Unable to get token: " + err.Error())
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Unable to get token: " + err.Error())
	}

	if resp.StatusCode != 200 {
		var errorInfo struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}
		err := json.Unmarshal(data, &errorInfo)
		if err != nil {
			return "", errors.New("Unable to get token: " + resp.Status)
		}
		return "", errors.New("Unable to get token: " + resp.Status + "; " + errorInfo.Message)
	}

	var iamToken struct {
		IamToken string `json:"iamToken"`
	}
	err = json.Unmarshal(data, &iamToken)
	if err != nil {
		return "", errors.New("Unable to unmarshal respose: " + err.Error())
	}
	return strings.TrimSpace(iamToken.IamToken), nil
}
