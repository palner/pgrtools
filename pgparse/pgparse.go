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
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/google/uuid"
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
			log.Printf("checkfields: %s exists in map and has the value %v", key, mapstring[key])
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

func GetUUID() string {
	u := uuid.New()
	return fmt.Sprintf("%s", u.String())
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
