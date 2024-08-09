/*
 * Copyright (C) 2021, 2024	The Palner Group, Inc. (palner.com)
 *							Fred Posner (@fredposner)
 *
 * pgkamtools is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version
 *
 * pgkamtools is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301  USA
 *
 */

package pgkamtools

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type StructAoRParse struct {
	AoR          json.RawMessage `json:"aor"`
	Address      json.RawMessage `json:"address"`
	Expires      json.RawMessage `json:"expires"`
	UA           json.RawMessage `json:"user-agent"`
	LastModified string          `json:"last-modified"`
}

type StructHtableDump struct {
	Jsonrpc string          `json:"jsonrpc"`
	Id      json.RawMessage `json:"id"`
	Result  []struct {
		Entry any `json:"entry"`
		Size  any `json:"size"`
		Slot  []struct {
			Name  json.RawMessage `json:"name"`
			Value json.RawMessage `json:"value"`
			Type  string          `json:"type"`
		} `json:"slot"`
	} `json:"result"`
}

type StructNameValue struct {
	Name  json.RawMessage `json:"name"`
	Value json.RawMessage `json:"value"`
}

type StructName struct {
	Name json.RawMessage `json:"name"`
}

type StructUserDump struct {
	Jsonrpc json.RawMessage `json:"jsonrpc"`
	Id      json.RawMessage `json:"id"`
	Result  struct {
		Domains []struct {
			Domain struct {
				Domain json.RawMessage `json:"Domain"`
				Size   json.RawMessage `json:"Size"`
				AoRs   []struct {
					Info struct {
						AoR      json.RawMessage `json:"AoR"`
						HashID   json.RawMessage `json:"HashID"`
						Contacts []struct {
							Contact struct {
								Address       json.RawMessage `json:"Address,omitempty"`
								Expires       json.RawMessage `json:"Expires,omitempty"`
								Q             json.RawMessage `json:"Q,omitempty"`
								CallID        json.RawMessage `json:"Call-ID,omitempty"`
								CSeq          json.RawMessage `json:"CSeq,omitempty"`
								UserAgent     json.RawMessage `json:"User-Agent,omitempty"`
								Received      json.RawMessage `json:"Received,omitempty"`
								Path          json.RawMessage `json:"Path,omitempty"`
								Socket        json.RawMessage `json:"Socket,omitempty"`
								Methods       json.RawMessage `json:"Methods,omitempty"`
								Ruid          json.RawMessage `json:"Ruid,omitempty"`
								Instance      json.RawMessage `json:"Instance,omitempty"`
								RegID         json.RawMessage `json:"Reg-Id,omitempty"`
								ServerID      json.RawMessage `json:"Server-Id,omitempty"`
								TcpconID      json.RawMessage `json:"Tcpconn-Id,omitempty"`
								Keepalive     json.RawMessage `json:"Keepalive,omitempty"`
								LastKeepAlive json.RawMessage `json:"Last-Keepalive,omitempty"`
								KARoundtrip   json.RawMessage `json:"KA-Roundtrip,omitempty"`
								LastModified  json.RawMessage `json:"Last-Modified,omitempty"`
							} `json:"Contact"`
						} `json:"Contacts"`
					} `json:"Info"`
				} `json:"AoRs"`
				Stats struct {
					Records  json.RawMessage `json:"Records"`
					MaxSlots json.RawMessage `json:"Max-Slots"`
				} `json:"Stats"`
			} `json:"Domain"`
		} `json:"Domains"`
	} `json:"result"`
}

type StructValue struct {
	Value json.RawMessage `json:"value"`
}

func CheckFields(mapstring map[string]string, reqfields []string) (bool, error) {
	errstring := ""
	for _, key := range reqfields {
		if _, exists := mapstring[key]; !exists {
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

func DispatcherAdd(groupval string, addressval string, urlval string) (string, error) {
	sendJson, _ := sjson.Set("", "jsonrpc", "2.0")
	sendJson, _ = sjson.Set(sendJson, "method", "dispatcher.add")
	sendJson, _ = sjson.Set(sendJson, "params.group", groupval)
	sendJson, _ = sjson.Set(sendJson, "params.address", addressval)
	sendJson, _ = sjson.Set(sendJson, "id", getId())

	results, err := SendJsonhttp(sendJson, urlval)
	if err != nil {
		return "", err
	}

	return results, nil
}

func DispatcherList(urlval string) (string, error) {
	sendJson, _ := sjson.Set("", "jsonrpc", "2.0")
	sendJson, _ = sjson.Set(sendJson, "method", "dispatcher.list")
	sendJson, _ = sjson.Set(sendJson, "id", getId())

	results, err := SendJsonhttp(sendJson, urlval)
	if err != nil {
		return "", err
	}

	return results, nil

}

func DispatcherListSimple(urlval string) (string, error) {
	sendJson, _ := sjson.Set("", "jsonrpc", "2.0")
	sendJson, _ = sjson.Set(sendJson, "method", "dispatcher.list")
	sendJson, _ = sjson.Set(sendJson, "id", getId())

	results, err := SendJsonhttp(sendJson, urlval)
	if err != nil {
		return "", err
	}

	if !gjson.Valid(results) {
		return "", errors.New("invalid response from kamailio")
	}

	if gjson.Get(results, "error.message").Exists() {
		errstring := gjson.Get(results, "error.message")
		return "", errors.New(errstring.String())
	}

	resultJson := gjson.Get(results, "result.RECORDS.#[@flatten].SET.TARGETS.#.DEST.URI")
	var jsonResult string
	for _, nodeValue := range resultJson.Array() {
		jsonResult, _ = sjson.Set(jsonResult, "nodes.-1", nodeValue.Str)
	}

	return jsonResult, nil
}

func DispatcherListByGroup(urlval string) (string, error) {
	sendJson, _ := sjson.Set("", "jsonrpc", "2.0")
	sendJson, _ = sjson.Set(sendJson, "method", "dispatcher.list")
	sendJson, _ = sjson.Set(sendJson, "id", getId())

	results, err := SendJsonhttp(sendJson, urlval)
	if err != nil {
		return "", err
	}

	if !gjson.Valid(results) {
		return "", errors.New("invalid response from kamailio")
	}

	if gjson.Get(results, "error.message").Exists() {
		errstring := gjson.Get(results, "error.message")
		return "", errors.New(errstring.String())
	}

	resultJson := gjson.Get(results, "result.RECORDS.#.SET.{id:ID,nodes:TARGETS.#.{uri:DEST.URI,flags:DEST.FLAGS,priority:DEST.PRIORITY,latency:DEST.LATENCY.AVG}}")
	return resultJson.String(), nil
}

func DispatcherRemove(groupval string, addressval string, urlval string) (string, error) {
	sendJson, _ := sjson.Set("", "jsonrpc", "2.0")
	sendJson, _ = sjson.Set(sendJson, "method", "dispatcher.remove")
	sendJson, _ = sjson.Set(sendJson, "params.group", groupval)
	sendJson, _ = sjson.Set(sendJson, "params.address", addressval)
	sendJson, _ = sjson.Set(sendJson, "id", getId())

	results, err := SendJsonhttp(sendJson, urlval)
	if err != nil {
		return "", err
	}

	return results, nil
}

func HtableDelete(tableval string, keyval string, urlval string) (bool, error) {
	sendJsonStr, _ := sjson.Set("", "jsonrpc", "2.0")
	sendJsonStr, _ = sjson.Set(sendJsonStr, "method", "htable.delete")
	sendJsonStr, _ = sjson.Set(sendJsonStr, "params.htable", tableval)
	sendJsonStr, _ = sjson.Set(sendJsonStr, "params.key", keyval)
	sendJsonStr, _ = sjson.Set(sendJsonStr, "id", getId())
	_, err := SendJsonhttp(sendJsonStr, urlval)

	if err != nil {
		return false, err
	}

	return true, nil
}

func HtableDump(tableval string, urlval string) (string, error) {
	sendJsonStr, _ := sjson.Set("", "jsonrpc", "2.0")
	sendJsonStr, _ = sjson.Set(sendJsonStr, "method", "htable.dump")
	sendJsonStr, _ = sjson.Set(sendJsonStr, "params.htable", tableval)
	sendJsonStr, _ = sjson.Set(sendJsonStr, "id", getId())
	htableresult, err := SendJsonhttp(sendJsonStr, urlval)

	if err != nil {
		return "", err
	}

	return htableresult, nil
}

func HtableFlush(tableval string, urlval string) (bool, error) {
	sendJsonStr, _ := sjson.Set("", "jsonrpc", "2.0")
	sendJsonStr, _ = sjson.Set(sendJsonStr, "method", "htable.flush")
	sendJsonStr, _ = sjson.Set(sendJsonStr, "params.htable", tableval)
	sendJsonStr, _ = sjson.Set(sendJsonStr, "id", getId())
	_, err := SendJsonhttp(sendJsonStr, urlval)

	if err != nil {
		return false, err
	}

	return true, nil
}

func HtableGet(tableval string, keyval string, urlval string) (string, error) {
	sendJsonStr, _ := sjson.Set("", "jsonrpc", "2.0")
	sendJsonStr, _ = sjson.Set(sendJsonStr, "method", "htable.get")
	sendJsonStr, _ = sjson.Set(sendJsonStr, "params.htable", tableval)
	sendJsonStr, _ = sjson.Set(sendJsonStr, "params.key", keyval)
	sendJsonStr, _ = sjson.Set(sendJsonStr, "id", getId())
	getval, err := SendJsonhttp(sendJsonStr, urlval)

	if err != nil {
		return "", err
	}

	parse, err := HtableParseValueSingle(getval)
	if err != nil {
		return "", err
	}

	return parse, nil
}

// changed 2023-01-18 to treat string as int in json for seti.
func HtableSetInt(tableval string, keyval string, valval string, urlval string) (bool, error) {
	sendJsonStr, _ := sjson.Set("", "jsonrpc", "2.0")
	sendJsonStr, _ = sjson.Set(sendJsonStr, "method", "htable.seti")
	sendJsonStr, _ = sjson.Set(sendJsonStr, "params.htable", tableval)
	sendJsonStr, _ = sjson.Set(sendJsonStr, "params.key", keyval)
	sendJsonStr, _ = sjson.Set(sendJsonStr, "params.value", valval)
	sendJsonStr, _ = sjson.Set(sendJsonStr, "id", getId())
	_, err := SendJsonhttp(sendJsonStr, urlval)

	if err != nil {
		return false, err
	}

	return true, nil
}

func HtableSetString(tableval string, keyval string, valval string, urlval string) (bool, error) {
	sendJsonStr, _ := sjson.Set("", "jsonrpc", "2.0")
	sendJsonStr, _ = sjson.Set(sendJsonStr, "method", "htable.sets")
	sendJsonStr, _ = sjson.Set(sendJsonStr, "params.htable", tableval)
	sendJsonStr, _ = sjson.Set(sendJsonStr, "params.key", keyval)
	sendJsonStr, _ = sjson.Set(sendJsonStr, "params.value", valval)
	sendJsonStr, _ = sjson.Set(sendJsonStr, "id", getId())
	_, err := SendJsonhttp(sendJsonStr, urlval)

	if err != nil {
		return false, err
	}

	return true, nil
}

func HtableParseNameOnly(jsonval string) (string, error) {
	if !gjson.Valid(jsonval) {
		return "", errors.New("invalid json")
	}

	if gjson.Get(jsonval, "error.message").Exists() {
		errstring := gjson.Get(jsonval, "error.message")
		return "", errors.New(errstring.String())
	}

	var dump StructHtableDump
	err := json.Unmarshal([]byte(jsonval), &dump)
	if err != nil {
		return "", err
	}

	tempparse := StructName{}
	var parsedNameVal []StructName
	for _, u := range dump.Result {
		for _, s := range u.Slot {
			tempparse.Name = s.Value
			parsedNameVal = append(parsedNameVal, tempparse)
		}
	}

	jsonByteData, err := json.MarshalIndent(parsedNameVal, "", "\t")
	if err != nil {
		jsonstr, _ := sjson.Set("", "error", true)
		jsonstr, _ = sjson.Set(jsonstr, "details", err.Error())
		return jsonstr, err
	}

	jsonStringData := string(jsonByteData)
	return jsonStringData, nil
}

func HtableParseNameValue(jsonval string) (string, error) {
	if !gjson.Valid(jsonval) {
		return "", errors.New("invalid json received")
	}

	if gjson.Get(jsonval, "error.message").Exists() {
		errstring := gjson.Get(jsonval, "error.message")
		return "", errors.New(errstring.String())
	}

	var dump StructHtableDump
	err := json.Unmarshal([]byte(jsonval), &dump)
	if err != nil {
		return "", err
	}

	tempparse := StructNameValue{}
	var parsedNameVal []StructNameValue
	for _, u := range dump.Result {
		for _, s := range u.Slot {
			tempparse.Name = s.Name
			tempparse.Value = s.Value
			parsedNameVal = append(parsedNameVal, tempparse)
		}
	}

	jsonByteData, err := json.MarshalIndent(parsedNameVal, "", "\t")
	if err != nil {
		jsonstr, _ := sjson.Set("", "error", true)
		jsonstr, _ = sjson.Set(jsonstr, "details", err.Error())
		return jsonstr, err
	}

	jsonStringData := string(jsonByteData)
	return jsonStringData, nil
}

func HtableParseValueOnly(jsonval string) (string, error) {
	if !gjson.Valid(jsonval) {
		return "", errors.New("invalid json received")
	}

	if gjson.Get(jsonval, "error.message").Exists() {
		errstring := gjson.Get(jsonval, "error.message")
		return "", errors.New(errstring.String())
	}

	var dump StructHtableDump
	err := json.Unmarshal([]byte(jsonval), &dump)
	if err != nil {
		return "", err
	}

	tempparse := StructValue{}
	var parsedNameVal []StructValue
	for _, u := range dump.Result {
		for _, s := range u.Slot {
			tempparse.Value = s.Value
			parsedNameVal = append(parsedNameVal, tempparse)
		}
	}

	jsonByteData, err := json.MarshalIndent(parsedNameVal, "", "\t")
	if err != nil {
		jsonstr, _ := sjson.Set("", "error", true)
		jsonstr, _ = sjson.Set(jsonstr, "details", err.Error())
		return jsonstr, err
	}

	jsonStringData := string(jsonByteData)
	return jsonStringData, nil
}

func HtableParseValueSingle(jsonval string) (string, error) {
	if !gjson.Valid(jsonval) {
		return "", errors.New("invalid json")
	}

	if gjson.Get(jsonval, "error.message").Exists() {
		errstring := gjson.Get(jsonval, "error.message")
		return "", errors.New(errstring.String())
	}

	parsedval := gjson.Get(jsonval, "result.item.{value:value}")
	return parsedval.String(), nil
}

func RegDeleteAOR(aorval string, urlval string) (bool, error) {
	sendJsonStr, _ := sjson.Set("", "jsonrpc", "2.0")
	sendJsonStr, _ = sjson.Set(sendJsonStr, "method", "ul.rm")
	sendJsonStr, _ = sjson.Set(sendJsonStr, "params.table", "location")
	sendJsonStr, _ = sjson.Set(sendJsonStr, "params.AOR", aorval)
	sendJsonStr, _ = sjson.Set(sendJsonStr, "id", getId())
	_, err := SendJsonhttp(sendJsonStr, urlval)

	if err != nil {
		return false, err
	}

	return true, nil
}

func RegGetAOR(aorval string, urlval string) (string, error) {
	sendJsonStr, _ := sjson.Set("", "jsonrpc", "2.0")
	sendJsonStr, _ = sjson.Set(sendJsonStr, "method", "ul.lookup")
	sendJsonStr, _ = sjson.Set(sendJsonStr, "params.table", "location")
	sendJsonStr, _ = sjson.Set(sendJsonStr, "params.AOR", aorval)
	sendJsonStr, _ = sjson.Set(sendJsonStr, "id", getId())
	aorresult, err := SendJsonhttp(sendJsonStr, urlval)

	if err != nil {
		return "", err
	}

	parsed, err := RegAorParse(aorresult)

	if err != nil {
		return "", err
	}

	return parsed, nil
}

func RegAorParse(jsonval string) (string, error) {
	if !gjson.Valid(jsonval) {
		return "", errors.New("invalid json")
	}

	if gjson.Get(jsonval, "error.message").Exists() {
		errstring := gjson.Get(jsonval, "error.message")
		return "", errors.New(errstring.String())
	}

	var dump StructUserDump
	err := json.Unmarshal([]byte(jsonval), &dump)
	if err != nil {
		return "", err
	}

	userparse := StructAoRParse{}
	var parsedUsers []StructAoRParse
	for _, d := range dump.Result.Domains {
		for _, a := range d.Domain.AoRs {
			userparse.AoR = a.Info.AoR
			for _, c := range a.Info.Contacts {
				userparse.Address = c.Contact.Address
				userparse.Expires = c.Contact.Expires
				userparse.UA = c.Contact.UserAgent
				lastmodified, _ := formatLastModifed(string(c.Contact.LastModified[:]))
				userparse.LastModified = lastmodified
				parsedUsers = append(parsedUsers, userparse)
			}
		}
	}

	jsonByteData, err := json.MarshalIndent(parsedUsers, "", "\t")
	if err != nil {
		return "", err
	}

	jsonStringData := string(jsonByteData)
	return jsonStringData, nil
}

func RegsAors(jsonval string) (string, error) {
	if !gjson.Valid(jsonval) {
		return "", errors.New("invalid json")
	}

	if gjson.Get(jsonval, "error.message").Exists() {
		errstring := gjson.Get(jsonval, "error.message")
		return "", errors.New(errstring.String())
	}

	parsedval := gjson.Get(jsonval, "result.Domains.#[@flatten].Domain.AoRs.#.Info.AoR")
	return parsedval.String(), nil
}

func RegsFullContactInfo(jsonval string) (string, error) {
	if !gjson.Valid(jsonval) {
		return "", errors.New("invalid json")
	}

	if gjson.Get(jsonval, "error.message").Exists() {
		errstring := gjson.Get(jsonval, "error.message")
		return "", errors.New(errstring.String())
	}

	parsedval := gjson.Get(jsonval, "result.Domains.#[@flatten].Domain.AoRs.#.{Info.AoR,Info.Contacts}.@ugly")
	return parsedval.String(), nil
}

func RegsGet(urlval string) (string, error) {
	sendjson := `{"jsonrpc": "2.0", "method": "ul.dump", "id":` + getId() + `}`
	htableresult, err := SendJsonhttp(sendjson, urlval)

	if err != nil {
		return "", err
	}

	return htableresult, nil
}

func RegsSimpleParse(jsonval string) (string, error) {
	if !gjson.Valid(jsonval) {
		return "", errors.New("invalid json")
	}

	if gjson.Get(jsonval, "error.message").Exists() {
		errstring := gjson.Get(jsonval, "error.message")
		return "", errors.New(errstring.String())
	}

	parsedval := gjson.Get(jsonval, "result.Domains.#[@flatten].Domain.AoRs.#.Info.{aor:AoR,details:Contacts.#.Contact.{address:Address,ua:User-Agent,expires:Expires,last-modified:Last-Modified}}")
	return parsedval.String(), nil
}

func RegsTotal(jsonval string) (string, error) {
	if !gjson.Valid(jsonval) {
		return "", errors.New("invalid json")
	}

	if gjson.Get(jsonval, "error.message").Exists() {
		errstring := gjson.Get(jsonval, "error.message")
		return "", errors.New(errstring.String())
	}

	parsedval := gjson.Get(jsonval, "result.Domains.#[@flatten].Domain.Stats.{Total_Registered:Records}")
	return parsedval.String(), nil
}

func RemoveDuplicatesUnordered(elements []string) []string {
	encountered := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	result := []string{}
	for key := range encountered {
		result = append(result, key)
	}

	return result
}

func SendJsonhttp(jsonstr string, urlstr string) (string, error) {
	var err error

	// send json to url
	sendbody := strings.NewReader(jsonstr)
	req, err := http.NewRequest("POST", urlstr, sendbody)

	if err != nil {
		// handle err
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	curlBody, err := io.ReadAll(resp.Body)

	if err != nil {
		// handle err
		return "error", err
	}

	return string(curlBody), nil
}

func SendJsonhttpIgnoreCert(jsonstr string, urlstr string) (string, error) {
	var err error
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr, Timeout: 2 * time.Second}

	// send json to url
	sendbody := strings.NewReader(jsonstr)
	req, err := http.NewRequest("POST", urlstr, sendbody)

	if err != nil {
		// handle err
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	curlBody, err := io.ReadAll(resp.Body)

	if err != nil {
		// handle err
		return "error", err
	}

	return string(curlBody), nil
}

// send a get request via http and return the response
func SendGethttp(urlstr string) (string, error) {
	var err error

	req, err := http.NewRequest("GET", urlstr, nil)

	if err != nil {
		// handle err
		return "error", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "error", err
	}

	defer resp.Body.Close()
	curlBody, err := io.ReadAll(resp.Body)

	if err != nil {
		// handle err
		return "error", err
	}

	return string(curlBody), nil
}

func SendGethttpIgnoreCert(urlstr string) (string, error) {
	var err error
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr, Timeout: 2 * time.Second}
	req, err := http.NewRequest("GET", urlstr, nil)
	if err != nil {
		// handle err
		return "error", err
	}

	resp, err := client.Do(req)
	if err != nil {
		if os.IsTimeout(err) {
			return "timeout", err
		}

		return "error", err
	}

	defer resp.Body.Close()
	curlBody, err := io.ReadAll(resp.Body)

	if err != nil {
		// handle err
		return "error", err
	}

	return string(curlBody), nil
}

func Uptime(urlval string) (string, error) {
	sendJsonStr, _ := sjson.Set("", "jsonrpc", "2.0")
	sendJsonStr, _ = sjson.Set(sendJsonStr, "method", "core.uptime")
	sendJsonStr, _ = sjson.Set(sendJsonStr, "id", getId())
	results, err := SendJsonhttp(sendJsonStr, urlval)

	if err != nil {
		return "", err
	}

	uptimeResult, _ := UptimeParse(results)
	return uptimeResult, nil
}

func UptimeParse(jsonval string) (string, error) {
	if !gjson.Valid(jsonval) {
		return "", errors.New("invalid json")
	}

	if gjson.Get(jsonval, "error.message").Exists() {
		errstring := gjson.Get(jsonval, "error.message")
		return "", errors.New(errstring.String())
	}

	parsedval := gjson.Get(jsonval, "result.uptime")
	return parsedval.String(), nil
}

func Version(urlval string) (string, error) {
	sendJsonStr, _ := sjson.Set("", "jsonrpc", "2.0")
	sendJsonStr, _ = sjson.Set(sendJsonStr, "method", "core.version")
	sendJsonStr, _ = sjson.Set(sendJsonStr, "id", getId())
	results, err := SendJsonhttp(sendJsonStr, urlval)

	if err != nil {
		return "", err
	}

	versionResult, _ := VersionParse(results)
	return versionResult, nil
}

func VersionParse(jsonval string) (string, error) {
	if !gjson.Valid(jsonval) {
		return "", errors.New("invalid json")
	}

	if gjson.Get(jsonval, "error.message").Exists() {
		errstring := gjson.Get(jsonval, "error.message")
		return "", errors.New(errstring.String())
	}

	parsedval := gjson.Get(jsonval, "result")
	return parsedval.String(), nil
}

func getId() string {
	timenow := time.Now().UnixMicro()
	timenowstr := strconv.FormatInt(timenow, 10)
	return timenowstr
}

func formatLastModifed(scientificNotation string) (string, error) {
	flt, _, err := big.ParseFloat(scientificNotation, 10, 0, big.ToNearestEven)
	if err != nil {
		return "", err
	}

	fltVal := fmt.Sprintf("%.0f", flt)
	intVal, err := strconv.ParseInt(fltVal, 10, 64)
	if err != nil {
		return "", err
	}

	lms := fmt.Sprint(uint(intVal))
	var lastmodified string
	i, err := strconv.ParseInt(lms, 10, 64)
	if err != nil {
		return "", err
	} else {
		tm := time.Unix(i, 0)
		lastmodified = fmt.Sprint(tm)
	}

	return lastmodified, nil
}
