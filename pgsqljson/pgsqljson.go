/*

Copyright (C) 2020, 2024 Fred Posner. All Rights Reserved.
Copyright (C) 2020, 2024 The Palner Group, Inc. All Rights Reserved.

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
