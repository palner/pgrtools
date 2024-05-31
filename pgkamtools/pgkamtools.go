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
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

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

	parsedval := gjson.Get(jsonval, "{\"result\":result.#.slot.#[@flatten].name}")
	return parsedval.String(), nil
}

func HtableParseNameValue(jsonval string) (string, error) {
	if !gjson.Valid(jsonval) {
		return "", errors.New("invalid json")
	}

	parsedval := gjson.Get(jsonval, "result.#.slot.#[@flatten].{name,value}.@pretty")
	return parsedval.String(), nil
}

func HtableParseValueOnly(jsonval string) (string, error) {
	if !gjson.Valid(jsonval) {
		return "", errors.New("invalid json")
	}

	parsedval := gjson.Get(jsonval, "{\"result\":result.#.slot.#[@flatten].value}")
	return parsedval.String(), nil
}

func HtableParseValueSingle(jsonval string) (string, error) {
	if !gjson.Valid(jsonval) {
		return "", errors.New("invalid json")
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

	parsedval := gjson.Get(jsonval, "result.Contacts.#[@flatten].{Contact.Address,Contact.Expires,Contact.User-Agent}")
	return parsedval.String(), nil
}

func RegsAors(jsonval string) (string, error) {
	if !gjson.Valid(jsonval) {
		return "", errors.New("invalid json")
	}

	parsedval := gjson.Get(jsonval, "result.Domains.#[@flatten].Domain.AoRs.#.Info.AoR")
	return parsedval.String(), nil
}

func RegsFullContactInfo(jsonval string) (string, error) {
	if !gjson.Valid(jsonval) {
		return "", errors.New("invalid json")
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

	parsedval := gjson.Get(jsonval, "result.Domains.#[@flatten].Domain.AoRs.#.Info.{aor:AoR,details:Contacts.#.Contact.{address:Address,ua:User-Agent,expires:Expires,last-modified:Last-Modified}}")
	return parsedval.String(), nil
}

func RegsTotal(jsonval string) (string, error) {
	if !gjson.Valid(jsonval) {
		return "", errors.New("invalid json")
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
	log.Print("sendJsonhttp request: ", jsonstr, " ", urlstr)
	var err error

	// send json to url
	sendbody := strings.NewReader(jsonstr)
	req, err := http.NewRequest("POST", urlstr, sendbody)

	if err != nil {
		// handle err
		log.Print(err)
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print(err)
		return "", err
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

// send a get request via http and return the response
func SendGethttp(urlstr string) (string, error) {
	var err error

	req, err := http.NewRequest("GET", urlstr, nil)

	if err != nil {
		// handle err
		log.Print(err)
		return "error", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print(err)
		return "error", err
	}

	defer resp.Body.Close()
	curlBody, err := io.ReadAll(resp.Body)

	if err != nil {
		// handle err
		log.Print(err)
		return "error", err
	}

	// log.Print("curl response -> ", string(curlBody))
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

	parsedval := gjson.Get(jsonval, "result")
	return parsedval.String(), nil
}

func getId() string {
	timenow := time.Now().UnixMicro()
	timenowstr := strconv.FormatInt(timenow, 10)
	return timenowstr
}
