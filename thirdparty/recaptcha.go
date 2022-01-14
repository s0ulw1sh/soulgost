package thirdparty

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const recaptchaServerPath = "https://www.google.com/recaptcha/api/siteverify"

type TRecaptchaResponse struct {
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

type TRecaptcha struct {
	Secret string
}

func (self *TRecaptcha) Check(response string) bool {
	var r TRecaptchaResponse

	resp, err := http.PostForm(recaptchaServerPath,
		url.Values{"secret": {self.Secret}, "response": {response}})
	
	if err != nil {
		return false
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return false
	}

	err = json.Unmarshal(body, &r)

	return err == nil && r.Success
}