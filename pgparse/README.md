# pgparse

pgparse is Part of [pgrtools](https://github.com/palner/pgrtools). See main site for license/warranty.

pgparse is designed to help with commonly needed parsing tools for web, api, json, etc.

## Functions

### Last10

Used commonly for VoIP and RTC operations, this function returns the last 10 characters of a string. If less than 10 digits, the string is returned.

#### Last10 Example

```go
newstring := pgparse.Last10("12345678910");
```

`newstring` would be: `2345678910`

### CheckFields

Checks to see if required fields are within a map[string]string.

#### CheckFields Example

```go
requiredfields := []string{"email"}
_, err = pgparse.CheckFields(keyVal, requiredfields)
if err != nil {
    log.Println(connectingip, "[getRecoverKey] parse error", err.Error())
    http.Redirect(w, r, "https://...", http.StatusMovedPermanently)
    return
}
```

### GetUUID

Returns a UUID value as a string

#### GetUUID Example

```go
myUUID := pgparse.GetUUID()
```

### LowerKeys

Transforms the keys from a map[sting]string to lower case.

#### LowerKeys Example

```go
bodyVal := make(amp[string]string)
\\ make your bodyVal...

\\ change bodyVal keys to lower case
bodyVal = pgparse.LowerKeys(bodyval)
```

### ParseBody

Parses an r.body `[]byte` and makes a map[string]string of the Json or QueryString objects.

#### ParseBody Example

```go
body, err := ioutil.ReadAll(r.Body)
if err != nil {
    ...
}

keyVal := pgparse.ParseBody(body)
```

### ParseBodyErr

Takes an r.body `[]byte` and returns a map[string]string, error of the json or query strings.

#### ParseBodyErr Example

```go
body, err := ioutil.ReadAll(r.Body)
if err != nil {
    ...
}

keyVal, err := pgparse.ParseBodyErr(body)
```

### ParseBodyFields

Receives a r `http.Request` and required fields `[]string`. Returns a map[string]string, error.

#### ParseBodyFields Example

```go
requiredfields := []string{"email", "name", "age"}
keyVal, err := pgparse.ParseBodyFields(r, requiredfields)
```

### PgParseForm

Receives a r `http.Request` and parses all form fields. Returns a map[string]string, error.

#### PgParseForm Example

```go
keyVal, err := pgparse.PgParseForm(r)
```

### PgParseFormFields

Receives a r `http.Request` and required fields `[]string`. Parses all form fields and checks for required fields. Returns a map[string]string, error.

#### PgParseFormFields Example

```go
requiredfields := []string{"email", "name", "age"}
keyVal, err := pgparse.PgParseFormFields(r, requiredfields)
```
