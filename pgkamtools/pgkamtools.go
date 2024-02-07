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

func HtableDelete(tableval string, keyval string, urlval string) (bool, error) {
	sendjson := `{"jsonrpc": "2.0", "method": "htable.delete", "params":{"htable":"` + tableval + `", "key":"` + keyval + `"}, "id":` + getId() + `}`
	_, err := SendJsonhttp(sendjson, urlval)

	if err != nil {
		return false, err
	}

	return true, nil
}

func HtableDump(tableval string, urlval string) (string, error) {
	sendjson := `{"jsonrpc": "2.0", "method": "htable.dump", "params":{"name":"` + tableval + `"}, "id":` + getId() + `}`
	htableresult, err := SendJsonhttp(sendjson, urlval)

	if err != nil {
		return "", err
	}

	return htableresult, nil
}

func HtableFlush(tableval string, urlval string) (bool, error) {
	sendjson := `{"jsonrpc": "2.0", "method": "htable.flush", "params":{"htable":"` + tableval + `"}, "id":` + getId() + `}`
	_, err := SendJsonhttp(sendjson, urlval)

	if err != nil {
		return false, err
	}

	return true, nil
}

func HtableGet(tableval string, keyval string, urlval string) (string, error) {
	sendjson := `{"jsonrpc": "2.0", "method": "htable.get", "params":{"htable":"` + tableval + `", "key":"` + keyval + `"}, "id":` + getId() + `}`
	getval, err := SendJsonhttp(sendjson, urlval)
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
	sendjson := `{"jsonrpc": "2.0", "method": "htable.seti", "params":{"htable":"` + tableval + `", "key":"` + keyval + `", "value":` + valval + `}, "id":` + getId() + `}`
	_, err := SendJsonhttp(sendjson, urlval)

	if err != nil {
		return false, err
	}

	return true, nil
}

func HtableSetString(tableval string, keyval string, valval string, urlval string) (bool, error) {
	sendjson := `{"jsonrpc": "2.0", "method": "htable.sets", "params":{"htable":"` + tableval + `", "key":"` + keyval + `", "value":"` + valval + `"}, "id":` + getId() + `}`
	_, err := SendJsonhttp(sendjson, urlval)

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
	sendjson := `{"jsonrpc": "2.0", "method": "ul.rm", "params":{"table":"location", "AOR":"` + aorval + `"}, "id":` + getId() + `}`
	_, err := SendJsonhttp(sendjson, urlval)

	if err != nil {
		return false, err
	}

	return true, nil
}

func RegGetAOR(aorval string, urlval string) (string, error) {
	sendjson := `{"jsonrpc": "2.0", "method": "ul.lookup", "params":{"table":"location", "AOR":"` + aorval + `"}, "id":` + getId() + `}`
	aorresult, err := SendJsonhttp(sendjson, urlval)

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
	for key, _ := range encountered {
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

func getId() string {
	timenow := time.Now().UnixMicro()
	timenowstr := strconv.FormatInt(timenow, 10)
	return timenowstr
}
