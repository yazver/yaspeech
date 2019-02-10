package yaspeech

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

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
