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
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/jaevor/go-nanoid"
)

func Last10(val string) string {
	log.Print("last10: checking ", val)
	if len(val) >= 10 {
		processedString := val[len(val)-10:]
		log.Println("last10: returning", processedString)
		return processedString
	} else {
		log.Print("last10:", val, "is less than 10 digits")
		return val
	}
}

func CheckFields(mapstring map[string]string, reqfields []string) (bool, error) {
	errstring := ""

	for _, key := range reqfields {
		if _, exists := mapstring[key]; exists {
			if mapstring[key] != "" {
				//log.Printf("checkfields: %s exists in map and has the value %v", key, mapstring[key])
			} else {
				log.Printf("checkfields: %s exists but has no value %v", key, mapstring[key])
				errstring += key + " is missing. "
			}
		} else {
			log.Printf("checkfields: %s is not found", key)
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
	log.Println("ParseBody: body received ->", string(body))
	if json.Valid(body) {
		log.Println("ParseBody: body is json")
		json.Unmarshal(body, &bodyVal)
		log.Println("ParseBody: body unmarshalled:", bodyVal)
	} else {
		log.Println("ParseBody: splitting based on &")
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
		log.Println("ParseBody: body split:", bodyVal)
	}

	return bodyVal
}

func ParseBodyErr(body []byte) (map[string]string, error) {
	bodyVal := make(map[string]string)
	log.Println("ParseBody: parsing body")
	if json.Valid(body) {
		log.Println("ParseBody: body is json")
		json.Unmarshal(body, &bodyVal)
	} else {
		if strings.Contains(string(body), "&") {
			log.Println("ParseBody: splitting based on &")
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
			log.Println("ParseBody: unable to parse body")
			return bodyVal, errors.New("unable to parse body. is it nil?")
		}
	}

	return bodyVal, nil
}

func ParseBodyFields(r *http.Request, reqfields []string) (map[string]string, error) {
	bodyVal := make(map[string]string)
	log.Println("ParseBodyFields: parsing body")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("ParseBodyFields: error reading body:", err.Error())
		return bodyVal, err
	}

	bodyVal, err = ParseBodyErr(body)
	if err != nil {
		log.Println("ParseBodyFields: error parsing body:", err.Error())
		return bodyVal, err
	}

	_, err = CheckFields(bodyVal, reqfields)
	if err != nil {
		log.Println("ParseBodyFields: error parsing body:", err.Error())
		return bodyVal, err
	}

	return bodyVal, nil
}

func PgParseForm(r *http.Request) (map[string]string, error) {
	bodyVal := make(map[string]string)
	err := r.ParseForm()
	if err != nil {
		log.Println("PgParseForm: error received -", err)
		return bodyVal, err
	}

	for key := range r.Form {
		log.Println("PgParseForm:", key, r.FormValue(key))
		bodyVal[key] = r.FormValue(key)
	}

	return bodyVal, nil
}

func PgParseFormFields(r *http.Request, reqfields []string) (map[string]string, error) {
	bodyVal := make(map[string]string)
	err := r.ParseForm()
	if err != nil {
		log.Println("PgParseFormFields: error received -", err)
		return bodyVal, err
	}

	for key := range r.Form {
		log.Println("PgParseFormFields:", key, r.FormValue(key))
		bodyVal[key] = r.FormValue(key)
	}

	_, err = CheckFields(bodyVal, reqfields)
	if err != nil {
		log.Println("PgParseFormFields: error parsing body:", err.Error())
		return bodyVal, err
	}

	return bodyVal, nil
}
