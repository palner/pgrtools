/*
 * Copyright (C) 2020 Fred Posner (The Palner Group, Inc.) (palner.com)
 *
 * This file is part of pgrtools, free software.
 *
 * pgrtools is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version
 *
 * pgrgotools is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301  USA
 *
 */

package pgparse

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/jaevor/go-nanoid"
)

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
