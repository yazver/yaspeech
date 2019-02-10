package yaspeech

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
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
			Code    string `json:"code"`
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

// The Voice for the synthesized speech.
// You can choose one of the following voices:
// Female voice: alyss, jane, oksana and omazh.
// Male voice: zahar and ermil.
// Default value of the parameter: oksana.
type Voice string

// Voices of the synthesized speech.
const (
	VoiceAlyss  Voice = "alyss"
	VoiceJane   Voice = "jane"
	VoiceOksana Voice = "oksana"
	VoiceOmazh  Voice = "omazh"
	VoiceZahar  Voice = "zahar"
	VoiceErmil  Voice = "ermil"
)

// Emotion is emotional tone of the voice.
// Acceptable values:
// good — Cheerful and friendly.
// evil — Irritated.
// neutral (default) — Without emotion.
type Emotion string

//Emotional tone of the voice
const (
	EmotionGood    Emotion = "good"
	EmotionEvil    Emotion = "evil"
	EmotionNeutral Emotion = "neutral"
)

// The Format of the synthesized audio.
// Acceptable values:
// lpcm — Audio file is synthesized in the LPCM format with no WAV header. Audio characteristics:
// Sampling — 8, 16, or 48 kHz, depending on the sampleRateHertz parameter value.
// Bit depth — 16-bit.
// Byte order — Reversed (little-endian).
// Audio data is stored as signed integers.
// oggopus (default) — Data in the audio file is encoded using the OPUS audio codec and compressed using the OGG container format (OggOpus).
type Format string

// Formats of the synthesized audio.
const (
	FormatLpcm    Format = "lpcm"
	FormatOggOpus Format = "oggopus"
)

// Language of the synthesized speech
// Acceptable values:
// ru-RU (default) — Russian.
// en-US — English.
// tr-TR — Turkish.
type Language string

// Languages
const (
	LangRU Language = "ru-RU"
	LangEN Language = "en-US"
	LangTR Language = "tr-TR"
)

// SampleRateHertz is the sampling frequency of the synthesized audio.
// Used if format is set to lpcm. Acceptable values:
// 48000 (default) — Sampling rate of 48 kHz.
// 16000 — Sampling rate of 16 kHz.
// 8000 — Sampling rate of 8 kHz.1
type SampleRateHertz int

// The sampling frequency of the synthesized audio
const (
	SampleRate48000 SampleRateHertz = 48000
	SampleRate16000 SampleRateHertz = 16000
	SampleRate8000  SampleRateHertz = 8000
)

// TextToSpeech is
type TextToSpeech struct {
	Voice      Voice
	Emotion    Emotion
	FolderID   string
	Format     Format
	Lang       Language
	SampleRate SampleRateHertz
	Speed      float32 // The speech rate is set as a decimal number in the range from 0.1 to 3.0
}

// NewTextToSpeech create and initializes TextToSpeech
func NewTextToSpeech(folderID string) *TextToSpeech {
	tts := &TextToSpeech{
		Voice:      VoiceAlyss,
		Emotion:    EmotionNeutral,
		Format:     FormatOggOpus,
		Lang:       LangRU,
		SampleRate: SampleRate48000,
		Speed:      1.0,
		FolderID:   folderID,
	}
	return tts
}

// Synthesize convert text to speech
func (tts *TextToSpeech) Synthesize(text string) ([]byte, error) {
	if strings.TrimSpace(text) == "" {
		return nil, errors.New("Text must be not empty")
	}

	iamtoken, err := token.Get()
	if err != nil {
		return nil, err
	}

	data := url.Values{
		"text":            {text},
		"voice":           {string(tts.Voice)},
		"emotion":         {string(tts.Emotion)},
		"folderId":        {tts.FolderID},
		"format":          {string(tts.Format)},
		"speed":           {strconv.FormatFloat(math.Min(math.Max(float64(tts.Speed), 0.1), 3.0), 'f', 1, 32)},
		"lang":            {string(tts.Lang)},
		"sampleRateHertz": {strconv.FormatInt(int64(tts.SampleRate), 10)},
	}
	req, err := http.NewRequest("POST", "https://tts.api.cloud.yandex.net/speech/v1/tts:synthesize", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+iamtoken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	err = checkResponse(resp)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	audio, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return audio, nil
}

func checkResponse(resp *http.Response) error {
	if resp.StatusCode != 200 {
		var errorInfo struct {
			ErrorCode    string `json:"error_code"`
			ErrorMessage string `json:"error_message"`
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			err := json.Unmarshal(data, &errorInfo)
			if err != nil {
				return errors.New("Request failed: " + resp.Status)
			}
			return errors.New("Request failed: " + resp.Status + "; " + errorInfo.ErrorMessage)
		}
		return errors.New("Request failed: " + resp.Status)
	}
	return nil
}
