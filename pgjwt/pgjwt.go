/*

Copyright (C) 2024 Fred Posner. All Rights Reserved.
Copyright (C) 2024 The Palner Group, Inc. All Rights Reserved.

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

package pgjwt

import (
	"bytes"
	"crypto/tls"
	"errors"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/tidwall/sjson"
)

func CheckAuth(r *http.Request, keys map[string]string, jwtKeystring string) (string, error) {
	var token string
	var err error

	if _, exists := keys["token"]; exists {
		token = keys["token"]
	} else {
		token, err = CheckBearer(r)
		if err != nil {
			return "", err
		}
	}

	_, err = CheckToken(token, []byte(jwtKeystring))
	if err != nil {
		return "", err
	}

	return "ok", nil
}

func CheckBearer(r *http.Request) (string, error) {
	bearerToken := r.Header.Get("Authorization")
	if !strings.Contains(bearerToken, " ") {
		tokenerror := errors.New("no valid header found")
		return "", tokenerror
	}

	reqToken := strings.Split(bearerToken, " ")[1]
	if len(reqToken) < 10 {
		tokenerror := errors.New("no token found")
		return "", tokenerror
	} else {
		return reqToken, nil
	}
}

func CheckBearerToken(r *http.Request, jwtKeystring string) (string, error) {
	bearerToken := r.Header.Get("Authorization")
	if !strings.Contains(bearerToken, " ") {
		return "", errors.New("no valid header found")
	}

	reqToken := strings.Split(bearerToken, " ")[1]
	if len(reqToken) < 10 {
		return "", errors.New("no token found")
	}

	_, err := CheckToken(reqToken, []byte(jwtKeystring))
	if err != nil {
		return "", err
	}

	return "ok", nil
}

func CheckGuiaccess(r *http.Request, cookiename string, jwtKeystring string) (string, error) {
	// We can obtain the session token from the requests cookies, which come with every request
	tknStr, err := CheckCookie(r, cookiename)
	if err != nil {
		return "", err
	}

	_, err = CheckToken(tknStr, []byte(jwtKeystring))
	if err != nil {
		return "", err
	}

	return "ok", nil
}

func CheckGuiaccessUsername(r *http.Request, cookiename string, jwtKeystring string) (string, error) {
	// We can obtain the session token from the requests cookies, which come with every request
	tknStr, err := CheckCookie(r, cookiename)
	if err != nil {
		return "", err
	}

	username, err := CheckToken(tknStr, []byte(jwtKeystring))
	if err != nil {
		return "", err
	}

	return username, nil
}

func CheckCookie(r *http.Request, name string) (string, error) {
	// Read the cookie as normal.
	c, err := r.Cookie(name)
	if err != nil {
		if err == http.ErrNoCookie {
			return "", errors.New("no cookie found")
		}

		// For any other type of error, return a bad request status
		return "", err
	}

	// Get the JWT string from the cookie
	return c.Value, nil
}

func GenerateToken(username string, jwtKeystring string, minutes time.Duration) (string, time.Time, error) {
	// get an expiration time of minutes minutes from now
	expirationTime := time.Now().Add(minutes * time.Minute)

	// Create the JWT claims, which includes the username and expiry time
	type Claims struct {
		Username string `json:"username"`
		jwt.RegisteredClaims
	}

	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString([]byte(jwtKeystring))

	if err != nil {
		return "fail", expirationTime, err
	}

	return tokenString, expirationTime, nil
}

func GenerateApitoken(username string, jwtKeystring string, days int) (string, error) {
	// get an expiration time days days from now
	expirationTime := time.Now().AddDate(0, 0, days)

	type Claims struct {
		Username string `json:"username"`
		jwt.RegisteredClaims
	}

	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString([]byte(jwtKeystring))
	if err != nil {
		return "fail", err
	}

	return tokenString, nil
}

func CheckToken(tokenstr string, jwtKey []byte) (string, error) {
	type Claims struct {
		Username string `json:"username"`
		jwt.RegisteredClaims
	}

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tokenstr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "fail", err
		}

		return "fail", err
	}

	if !tkn.Valid {
		return "fail", errors.New("token is not valid")
	}

	return claims.Username, nil
}

func TotpGenerate(issuer string, account string) (*otp.Key, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: account,
	})

	if err != nil {
		return nil, err
	}

	return key, nil
}

func TotpGenerateImgSecret(issuer string, account string) (image.Image, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: account,
	})

	if err != nil {
		return nil, "", err
	}

	image, err := key.Image(400, 400)
	if err != nil {
		return nil, "", err
	}

	return image, key.Secret(), nil
}

func TotpQRcode(key *otp.Key) (string, error) {
	img, err := key.Image(200, 200)
	if err != nil {
		return "", err
	}

	uuid := uuid.New()
	file, _ := os.Create("/tmp/" + uuid.String() + ".png")
	defer file.Close()
	png.Encode(file, img)
	return "/tmp/" + uuid.String() + ".png", nil
}

func TotpQRcodeExport(key *otp.Key) (*bytes.Buffer, error) {
	img, err := key.Image(200, 200)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func TotpQRcodeImage(key *otp.Key) (image.Image, error) {
	img, err := key.Image(200, 200)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func DisplayQRCode(image image.Image, w http.ResponseWriter, r *http.Request) error {
	// ... (create and draw on 'img' as shown above) ...
	w.Header().Set("Content-Type", "image/png")
	err := png.Encode(w, image)
	if err != nil {
		return err
	}

	return nil
}

func LastString(ss []string) string {
	return ss[len(ss)-1]
}

func FirstString(ss []string) string {
	return ss[0]
}

func HttpGet(urlstr string, seconds time.Duration) (string, error) {
	var err error
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   seconds * time.Second,
	}

	req, err := http.NewRequest("GET", urlstr, nil)
	if err != nil {
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
		return "error", err
	}

	return string(curlBody), nil
}

func JwtParseGetCert(tokenString string) (jwt.MapClaims, string, error) {
	// Parse and validate the JWT token
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, "error", err
	}

	x5u, ok := token.Header["x5u"].(string)
	if !ok || x5u == "" {
		return nil, "error", errors.New("no x5u header found")
	}

	cert, err := HttpGet(x5u, 2)
	if err != nil {
		return nil, "error", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, cert, nil
	} else {
		return nil, cert, err
	}
}

func JSONHandleError(w http.ResponseWriter, r *http.Request, errCode string, errDesc string, httpCode int64) {
	w.Header().Set("Content-Type", "application/json")
	httpJson, _ := sjson.Set("", "response", httpCode)
	httpJson, _ = sjson.Set(httpJson, "error", errCode)
	httpJson, _ = sjson.Set(httpJson, "details", errDesc)
	switch httpCode {
	case 400:
		w.WriteHeader(http.StatusBadRequest)
	case 403:
		w.WriteHeader(http.StatusForbidden)
	case 404:
		w.WriteHeader(http.StatusNotFound)
	case 408:
		w.WriteHeader(http.StatusRequestTimeout)
	case 429:
		w.WriteHeader(http.StatusTooManyRequests)
	case 500:
		w.WriteHeader(http.StatusInternalServerError)
	case 503:
		w.WriteHeader(http.StatusServiceUnavailable)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

	io.WriteString(w, httpJson+"\n")
}
