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

package pgsqljson

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

func CheckFields(mapstring map[string]string, reqfields []string) string {
	errstring := ""
	for _, key := range reqfields {
		if _, exists := mapstring[key]; !exists {
			errstring += key + " is missing. "
		}
	}

	return errstring
}

func ProcessResults(rows *sql.Rows) (string, error) {
	var err error
	cols, _ := rows.Columns()
	list := make([]map[string]interface{}, 0)
	for rows.Next() {
		vals := make([]interface{}, len(cols))
		for i := range cols {
			var s string
			vals[i] = &s
		}

		err = rows.Scan(vals...)
		if err != nil {
			log.Fatal(err)
		}

		m := make(map[string]interface{})
		for i, val := range vals {
			m[cols[i]] = val
		}

		list = append(list, m)
	}

	b, _ := json.MarshalIndent(list, "", "\t")
	jsonString := string(b)

	if err != nil {
		return "error", err
	}

	return jsonString, nil
}

func SendJsonhttp(jsonstr string, urlstr string) (string, error) {
	var err error

	// send json to url
	sendbody := strings.NewReader(jsonstr)
	req, err := http.NewRequest("POST", urlstr, sendbody)

	if err != nil {
		// handle err
		log.Print(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print(err)
	}

	defer resp.Body.Close()
	curlBody, err := io.ReadAll(resp.Body)

	if err != nil {
		// handle err
		log.Print(err)
		return "error", err
	}

	return string(curlBody), nil
}
