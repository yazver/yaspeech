package auth

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

const iamTokenLifetime = 12

// TokenGetter is the interface that wraps the GetToken method.
type TokenGetter interface {
	GetToken() (string, error)
}

// IAMToken is IAM token used for requests
type IAMToken struct {
	tokenGetter TokenGetter
	err         error
	token       string
	lastUpdate  time.Time
	mutex       sync.RWMutex
	ctx         context.Context
}

// NewIAMToken initialize IAM token.
func NewIAMToken(tokenGetter TokenGetter) *IAMToken {
	return NewIAMTokenContext(nil, tokenGetter)
}

// NewIAMTokenContext initialize IAM token.
// The token will be updated in the background.
func NewIAMTokenContext(ctx context.Context, tokenGetter TokenGetter) *IAMToken {
	t := &IAMToken{ctx: ctx, tokenGetter: tokenGetter}
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
	t.mutex.Lock()
	defer t.mutex.Unlock()

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

	if t.token != "" && time.Since(t.lastUpdate) < time.Hour*(iamTokenLifetime-4) {
		return
	}

	token, err := t.tokenGetter.GetToken()
	if err != nil && t.token != "" && time.Since(t.lastUpdate) < (time.Hour*(iamTokenLifetime-1)+time.Minute*50) {
		return
	}
	t.token, t.err, t.lastUpdate = token, err, time.Now()
}

func processResponse(r *http.Response, result interface{}) error {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusOK {
		var errorInfo struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}
		err := json.Unmarshal(data, &errorInfo)
		if err != nil {
			return errors.New(r.Status + ", " + string(data))
		}
		return errors.New(r.Status + "; " + errorInfo.Message)
	}

	err = json.Unmarshal(data, result)
	if err != nil {
		return errors.New("Unable to unmarshal respose: " + err.Error())
	}

	return nil
}

func processIAMResponse(r *http.Response) (string, error) {
	var iamToken struct {
		IamToken string `json:"iamToken"`
	}
	err := processResponse(r, &iamToken)
	if err != nil {
		return "", errors.New("Unable to get token: " + err.Error())
	}
	return strings.TrimSpace(iamToken.IamToken), nil
}

// ServiceAccounts list service accounts associated with the folder.
func ServiceAccounts(folderID string, token *IAMToken) (accounts map[string]string, err error) {
	accounts = make(map[string]string)
	iamtoken, err := token.Get()
	if err != nil {
		return
	}

	req, err := http.NewRequest("GET", "https://iam.api.cloud.yandex.net/iam/v1/serviceAccounts?folderId="+folderID, nil)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", "Bearer "+iamtoken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	var serviceAccounts struct {
		ServiceAccounts []struct {
			ID        string    `json:"id"`
			FolderID  string    `json:"folderId"`
			CreatedAt time.Time `json:"createdAt"`
			Name      string    `json:"name"`
		} `json:"serviceAccounts"`
	}
	err = processResponse(resp, &serviceAccounts)
	if err != nil {
		return
	}

	for _, s := range serviceAccounts.ServiceAccounts {
		accounts[s.Name] = s.ID
	}

	return
}

// CreateServiceAccountKey return the key ID and the private key used for initialization IAM token.
func CreateServiceAccountKey(id string, token *IAMToken) (keyID string, privateKey string, err error) {
	iamtoken, err := token.Get()
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "https://iam.api.cloud.yandex.net/iam/v1/keys", strings.NewReader(`{"serviceAccountId": "`+id+`"}`))
	if err != nil {
		return
	}
	req.Header.Add("Authorization", "Bearer "+iamtoken)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	var keyInfo struct {
		Key struct {
			CreatedAt        time.Time `json:"createdAt"`
			Description      string    `json:"description"`
			ID               string    `json:"id"`
			KeyAlgorithm     string    `json:"keyAlgorithm"`
			PublicKey        string    `json:"publicKey"`
			ServiceAccountID string    `json:"serviceAccountId"`
		} `json:"key"`
		PrivateKey string `json:"privateKey"`
	}
	err = processResponse(resp, &keyInfo)
	if err != nil {
		return
	}
	keyID = keyInfo.Key.ID
	privateKey = strings.Replace(keyInfo.PrivateKey, `\n`, "\n", -1)

	return
}
