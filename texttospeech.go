package yaspeech

import (
	"errors"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

	options := url.Values{
		"text":            {text},
		"voice":           {string(tts.Voice)},
		"emotion":         {string(tts.Emotion)},
		"folderId":        {tts.FolderID},
		"format":          {string(tts.Format)},
		"speed":           {strconv.FormatFloat(math.Min(math.Max(float64(tts.Speed), 0.1), 3.0), 'f', 1, 32)},
		"lang":            {string(tts.Lang)},
		"sampleRateHertz": {strconv.FormatInt(int64(tts.SampleRate), 10)},
	}
	req, err := http.NewRequest("POST", "https://tts.api.cloud.yandex.net/speech/v1/tts:synthesize", strings.NewReader(options.Encode()))
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
