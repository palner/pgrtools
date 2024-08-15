# pgparse

pgkamtools is Part of [pgrtools](https://github.com/palner/pgrtools). See main repo for license/warranty (tl;dr MIT).

pgkamtools is designed to help with commonly needed jsonrpc calls to kamailio.

## Docucmentation Welcome

Documentation / contributions gladly accepted.

## Background

pgkamtools expects a JSONRPC listener on Kamailio via http. Such as...

```
...
loadmodule "jsonrpcs.so"
loadmodule "xhttp.so"
...
modparam("jsonrpcs", "pretty_format", 1)
modparam("jsonrpcs", "transport", 1)
...
event_route[xhttp:request] {
	if ($hu =~ "^/RPC") {
		jsonrpc_dispatch();
		exit;
	} else {
		xhttp_reply("404", "NOT FOUND", "text/html", "$<html><body>404 - $hu [$si:$sp]</body></html>\n");
		exit;
	}

	return;
}
```

## Functions

### CheckFields

Expects

* map[string]string (values from form, etc)
* []string (list of required fields)

Returns

* bool
* error

#### Example

```go
body, err := ioutil.ReadAll(r.Body)
...
keyVal := pgparse.LowerKeys(pgparse.ParseBody(body))
...
_, err = pgkamtools.CheckFields(keyVal, []string{"group","address","url"})
...
```

(returns true if group, address, url are in keyVal)

### DispatcherAdd

See [Kamailio RPC dispatcher.add](https://kamailio.org/docs/modules/5.8.x/modules/dispatcher.html#dispatcher.r.add) documentation.

Expects

* group (string)
* address (string)
* url (string) (url for kamailio rpc)

Returns

* json string (result from kamailio)
* error

### DispatcherList

Returns all dispatcher info from kamailio

Expects

* url (string) (url for kamailio rpc)

### DipatcherListSimple

Returns simplified list of dispatcher nodes

Expects

* url (string) (url for kamailio rpc)

### DispatcherRemove

Expects

* group (string)
* address (string)
* url (string) (url for kamailio rpc)

Returns

* string
* error

### HtableDelete

Deletes a key from htables

Expects

* htable (string)
* key (strint)
* url (string) (url for kamailio rpc)

Returns

* bool
* error

Example

```go
_, err := pgkamtools.HtableDelete("ipban","1.1.1.1","http://localhost/RPC")
```

### HtableDump

Expects

* htable (string)
* url (string) (url for kamailio rpc)

Returns

* string (json)
* error

Example

```go
jsonData, err := pgkamtools.HtableDump("ipban","http://localhost/RPC")
```

### HtableFlush

### HtableGet

### HtableSetInt

### HtableSetString

### HtableParseNameOnly

### HtableParseNameValue

### HtableParseValueOnly

### HtableParseValueOnly

### RegDeleteAOR

### RegGetAOR

### RegAorParse

### RegsAors

### RegsFullContactInfo

### RegsGet

### RegsSimpleParse

### RegsTotal

### RemoveDuplicatesUnordered

### SendJsonhttp

### SendJsonhttpTimeout

### SendJsonhttpIgnoreCert

### SendJsonhttpIgnoreCertTimeout

### SendGethttp

### SendGethttpIgnoreCert

### SendGethttpIgnoreCertTimeout

### Uptime

### UptimeParse

### Version

### VersionParse

### getId

### formatLastModifed

