package pgkamtools

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

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

func GetAOR(aorval string, urlval string) (string, error) {
	sendjson := `{"jsonrpc": "2.0", "method": "ul.lookup", "params":{"table":"location", "AOR":"` + aorval + `"}, "id":1}`
	aorresult, err := SendJsonhttp(sendjson, urlval)

	if err != nil {
		return "", err
	}

	parsed, err := regAorParse(aorresult)

	if err != nil {
		return "", err
	}

	return parsed, nil
}

func GetHtable(tableval string, urlval string) (string, error) {
	sendjson := `{"jsonrpc": "2.0", "method": "htable.dump", "params":{"name":"` + tableval + `"}, "id":1}`
	htableresult, err := SendJsonhttp(sendjson, urlval)

	if err != nil {
		return "", err
	}

	return htableresult, nil
}

func GetRegs(urlval string) (string, error) {
	sendjson := `{"jsonrpc": "2.0", "method": "ul.dump",, "id":1}`
	htableresult, err := SendJsonhttp(sendjson, urlval)

	if err != nil {
		return "", err
	}

	return htableresult, nil
}

func HtableParseNameOnly(jsonval string) (string, error) {
	if !gjson.Valid(jsonval) {
		return "", errors.New("invalid json")
	}

	parsedval := gjson.Get(jsonval, "result.#.slot.#[@flatten].name")
	return parsedval.String(), nil
}

func HtableParseNameValue(jsonval string) (string, error) {
	if !gjson.Valid(jsonval) {
		return "", errors.New("invalid json")
	}

	parsedval := gjson.Get(jsonval, "result.#.slot.#[@flatten].{name,value}.@pretty")
	return parsedval.String(), nil
}

func RegsAors(jsonval string) (string, error) {
	if !gjson.Valid(jsonval) {
		return "", errors.New("invalid json")
	}

	parsedval := gjson.Get(jsonval, "result.Domains.#[@flatten].Domain.AoRs.#.Info.AoR")
	return parsedval.String(), nil
}

func RegAorParse(jsonval string) (string, error) {
	if !gjson.Valid(jsonval) {
		return "", errors.New("invalid json")
	}

	parsedval := gjson.Get(jsonval, "result.Contacts.#[@flatten].{Contact.Address,Contact.Expires,Contact.User-Agent}")
	return parsedval.String(), nil
}

func RegsFullContactInfo(jsonval string) (string, error) {
	if !gjson.Valid(jsonval) {
		return "", errors.New("invalid json")
	}

	parsedval := gjson.Get(jsonval, "result.Domains.#[@flatten].Domain.AoRs.#.{Info.AoR,Info.Contacts}.@ugly")
	return parsedval.String(), nil
}

func RegsSimpleParse(jsonval string) (string, error) {
	if !gjson.Valid(jsonval) {
		return "", errors.New("invalid json")
	}

	parsedval := gjson.Get(jsonval, "result.Domains.#[@flatten].Domain.AoRs.#.{Info.AoR,Info.Contacts.#[@flatten].Contact.Address,Info.Contacts.#[@flatten].Contact.Expires}")
	return parsedval.String(), nil
}

func RegsTotal(jsonval string) (string, error) {
	if !gjson.Valid(jsonval) {
		return "", errors.New("invalid json")
	}

	parsedval := gjson.Get(jsonval, "result.Domains.#[@flatten].Domain.Stats.Records")
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
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print(err)
	}

	defer resp.Body.Close()
	curlBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		// handle err
		log.Print(err)
		return "error", err
	}

	log.Print("curl response -> ", string(curlBody))
	return string(curlBody), nil
}
