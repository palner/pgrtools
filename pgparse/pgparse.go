/*

Copyright (C) 2020 Fred Posner. All Rights Reserved.
Copyright (C) 2020 The Palner Group, Inc. All Rights Reserved.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

package pgparse

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jaevor/go-nanoid"
	"github.com/tidwall/gjson"
)

type SendTurnstile struct {
	Secret   string `json:"secret"`
	Token    string `json:"response"`
	RemoteIp string `json:"remoteip"`
	Uuid     string `json:"idempotency_key"`
}

func Last10(val string) string {
	if len(val) >= 10 {
		processedString := val[len(val)-10:]
		return processedString
	} else {
		return val
	}
}

func CheckFields(mapstring map[string]string, reqfields []string) (bool, error) {
	errstring := ""
	for _, key := range reqfields {
		if _, exists := mapstring[key]; exists {
			if mapstring[key] != "" {
				// all good
			} else {
				errstring += key + " is missing. "
			}
		} else {
			errstring += key + " is missing. "
		}
	}

	if len(errstring) > 0 {
		err := errors.New(errstring)
		return false, err
	} else {
		return true, nil
	}
}

func CheckHoneypot(mapstring map[string]string, honeypot string) bool {
	if _, exists := mapstring[honeypot]; exists {
		if mapstring[honeypot] != "" {
			// honeypot is present
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func CheckFieldsAny(mapstring map[string]any, reqfields []string) (bool, error) {
	errstring := ""
	for _, key := range reqfields {
		if _, exists := mapstring[key]; exists {
			if mapstring[key] != "" {
				// all good
			} else {
				errstring += key + " is missing. "
			}
		} else {
			errstring += key + " is missing. "
		}
	}

	if len(errstring) > 0 {
		err := errors.New(errstring)
		return false, err
	} else {
		return true, nil
	}
}

func CheckTurnstile(mapstring map[string]string, secret string, remoteip string) (bool, error) {
	if _, exists := mapstring["cf-turnstile-response"]; exists {
		if mapstring["cf-turnstile-response"] == "" {
			return false, errors.New("cf-turnstile-reponse is empty")
		}
	} else {
		return false, errors.New("cf-turnstile-reponse is not present")
	}

	uuid, _ := getNanoID()

	// make paylod for cloudflare
	payload := SendTurnstile{
		Secret:   secret,
		Token:    mapstring["cf-turnstile-response"],
		RemoteIp: remoteip,
		Uuid:     uuid,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return false, err
	}

	// send json to url
	req, err := http.NewRequest("POST", "https://challenges.cloudflare.com/turnstile/v0/siteverify", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}

	req.Header = http.Header{
		"Content-Type": {"application/json"},
		"Accept":       {"application/json"},
		"User-Agent":   {"pgrtools"},
	}

	req.Header.Set("Content-Type", "application/json")
	http.DefaultClient.Timeout = 2 * time.Second
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	curlBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if !gjson.Valid(string(curlBody)) {
		return false, errors.New("invalid response from cloudlfare")
	}

	if gjson.Get(string(curlBody), "success").Exists() {
		val := gjson.Get(string(curlBody), "success")
		return val.Bool(), nil
	} else {
		return false, errors.New("success field not found in cloudflare response")
	}

}

func getNanoID() (string, error) {
	id, err := nanoid.Standard(21)
	if err != nil {
		return "", err
	} else {
		id1 := id()
		return id1, nil
	}

}

func getNanoIDSmall() (string, error) {
	id, err := nanoid.Standard(8)
	if err != nil {
		return "", err
	} else {
		id1 := id()
		return id1, nil
	}

}

func GetUUID() string {
	u := uuid.New()
	return u.String()
}

func LowerKeys(keyVal map[string]string) map[string]string {
	lf := make(map[string]string, len(keyVal))
	for k, v := range keyVal {
		lf[strings.ToLower(k)] = v
	}

	return lf
}

func ParseBody(body []byte) map[string]string {
	bodyVal := make(map[string]string)
	if json.Valid(body) {
		json.Unmarshal(body, &bodyVal)
	} else {
		stringsplit := strings.Split(string(body), "&")
		for _, pair := range stringsplit {
			z := strings.Split(pair, "=")
			decodedValue, err := url.QueryUnescape(z[1])
			if err != nil {
				bodyVal[z[0]] = z[1]
			} else {
				bodyVal[z[0]] = decodedValue
			}
		}
	}

	return bodyVal
}

func ParseBodyErr(body []byte) (map[string]string, error) {
	bodyVal := make(map[string]string)
	if json.Valid(body) {
		json.Unmarshal(body, &bodyVal)
	} else {
		if strings.Contains(string(body), "&") {
			stringsplit := strings.Split(string(body), "&")
			for _, pair := range stringsplit {
				z := strings.Split(pair, "=")
				decodedValue, err := url.QueryUnescape(z[1])
				if err != nil {
					bodyVal[z[0]] = z[1]
				} else {
					bodyVal[z[0]] = decodedValue
				}
			}
		} else {
			return bodyVal, errors.New("unable to parse body. is it nil?")
		}
	}

	return bodyVal, nil
}

func ParseBodyErrAny(body []byte) (map[string]any, error) {
	bodyVal := make(map[string]any)
	if json.Valid(body) {
		json.Unmarshal(body, &bodyVal)
	} else {
		if strings.Contains(string(body), "&") {
			stringsplit := strings.Split(string(body), "&")
			for _, pair := range stringsplit {
				z := strings.Split(pair, "=")
				decodedValue, err := url.QueryUnescape(z[1])
				if err != nil {
					bodyVal[z[0]] = z[1]
				} else {
					bodyVal[z[0]] = decodedValue
				}
			}
		} else {
			return bodyVal, errors.New("unable to parse body. is it nil?")
		}
	}

	return bodyVal, nil
}

func ParseBodyFields(r *http.Request, reqfields []string) (map[string]string, error) {
	bodyVal := make(map[string]string)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return bodyVal, err
	}

	bodyVal, err = ParseBodyErr(body)
	if err != nil {
		return bodyVal, err
	}

	_, err = CheckFields(bodyVal, reqfields)
	if err != nil {
		return bodyVal, err
	}

	return bodyVal, nil
}

func ParseBodyFieldsAny(r *http.Request, reqfields []string) (map[string]any, error) {
	bodyVal := make(map[string]any)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return bodyVal, err
	}

	bodyVal, err = ParseBodyErrAny(body)
	if err != nil {
		return bodyVal, err
	}

	_, err = CheckFieldsAny(bodyVal, reqfields)
	if err != nil {
		return bodyVal, err
	}

	return bodyVal, nil
}

func PgParseForm(r *http.Request) (map[string]string, error) {
	bodyVal := make(map[string]string)
	err := r.ParseForm()
	if err != nil {
		return bodyVal, err
	}

	for key := range r.Form {
		bodyVal[key] = r.FormValue(key)
	}

	return bodyVal, nil
}

func PgParseFormFields(r *http.Request, reqfields []string) (map[string]string, error) {
	bodyVal := make(map[string]string)
	err := r.ParseForm()
	if err != nil {
		return bodyVal, err
	}

	for key := range r.Form {
		bodyVal[key] = r.FormValue(key)
	}

	_, err = CheckFields(bodyVal, reqfields)
	if err != nil {
		return bodyVal, err
	}

	return bodyVal, nil
}
