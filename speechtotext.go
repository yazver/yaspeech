package yaspeech

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

// SpeechToText is
type SpeechToText struct {
	ProfanityFilter bool
	Topic           Topic
	FolderID        string
	Format          Format
	Lang            Language
	SampleRate      SampleRateHertz
	Speed           float32 // The speech rate is set as a decimal number in the range from 0.1 to 3.0
}

// NewSpeechToText create and initializes SpeechToText
func NewSpeechToText(folderID string) *SpeechToText {
	stt := &SpeechToText{
		ProfanityFilter: false,
		Topic:           TopicGeneral,
		Format:          FormatOggOpus,
		Lang:            LangRU,
		SampleRate:      SampleRate48000,
		Speed:           1.0,
		FolderID:        folderID,
	}
	return stt
}

// Recognize convert text to speech
func (stt *SpeechToText) Recognize(r io.Reader) (string, error) {
	iamtoken, err := token.Get()
	if err != nil {
		return "", err
	}

	options := url.Values{
		"topic":           {string(stt.Topic)},
		"profanityFilter": {strconv.FormatBool(stt.ProfanityFilter)},
		"folderId":        {stt.FolderID},
		"format":          {string(stt.Format)},
		"speed":           {strconv.FormatFloat(math.Min(math.Max(float64(stt.Speed), 0.1), 3.0), 'f', 1, 32)},
		"lang":            {string(stt.Lang)},
		"sampleRateHertz": {strconv.FormatInt(int64(stt.SampleRate), 10)},
	}
	req, err := http.NewRequest("POST", "https://stt.api.cloud.yandex.net/speech/v1/stt:recognize/?"+options.Encode(), r)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+iamtoken)
	req.Header.Add("Transfer-Encoding", "chunked")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	err = checkResponse(resp)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result struct {
		Result string `json:"result"`
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return "", errors.New("Unable to unmarshal respose: " + err.Error())
	}
	return result.Result, nil
}
