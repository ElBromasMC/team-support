package service

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type Captcha struct {
	siteKey   string
	secretKey string
}

func NewCaptchaService(siteKey, secretKey string) Captcha {
	return Captcha{
		siteKey:   siteKey,
		secretKey: secretKey,
	}
}

func (c Captcha) GetSiteKey() string {
	return c.siteKey
}

func (c Captcha) Verify(response, remoteIP string) (bool, float64, error) {
	verifyURL := "https://www.google.com/recaptcha/api/siteverify"

	data := url.Values{}
	data.Set("secret", c.secretKey)
	data.Set("response", response)
	data.Set("remoteip", remoteIP)

	resp, err := http.PostForm(verifyURL, data)
	if err != nil {
		return false, 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Success     bool     `json:"success"`
		Score       float64  `json:"score"` // This field is provided by reCAPTCHA v3.
		Action      string   `json:"action"`
		ChallengeTS string   `json:"challenge_ts"`
		Hostname    string   `json:"hostname"`
		ErrorCodes  []string `json:"error-codes"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, 0, err
	}

	return result.Success, result.Score, nil
}
