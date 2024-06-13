package pgjwt

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type SimpleJsonString struct {
	Json string
}

func checkAuth(r *http.Request, keys map[string]string, jwtKeystring string) (string, error) {
	connectingip := r.RemoteAddr
	log.Println(connectingip, "[checkAuth] request received")
	var token string
	var err error

	if _, exists := keys["token"]; exists {
		log.Println(connectingip, "[checkAuth] token found in body")
		token = keys["token"]
	} else {
		token, err = checkBearer(r)
		if err != nil {
			return "", err
		}
	}

	_, err = checkToken(token, []byte(jwtKeystring))
	if err != nil {
		return "", err
	}

	return "ok", nil
}

func checkBearer(r *http.Request) (string, error) {
	connectingip := r.RemoteAddr

	bearerToken := r.Header.Get("Authorization")
	if !strings.Contains(bearerToken, " ") {
		log.Println(connectingip, "[checkBearer] valid header not found")
		tokenerror := errors.New("no valid header found")
		return "", tokenerror
	}

	reqToken := strings.Split(bearerToken, " ")[1]
	if len(reqToken) < 10 {
		log.Println(connectingip, "[checkBearer] no token found")
		tokenerror := errors.New("no token found")
		return "", tokenerror
	} else {
		log.Println(connectingip, "[checkBearer] token found in bearer header")
		return reqToken, nil
	}
}

func checkGuiaccess(r *http.Request, cookiename string, jwtKeystring string) (string, error) {
	log.Print("checkGuiaccess: checking cookie", cookiename, "for token")

	// We can obtain the session token from the requests cookies, which come with every request
	tknStr, err := checkCookie(r, cookiename)
	if err != nil {
		log.Println("checkGuiaccess: cookie error", err.Error())
		return "", err
	}

	_, err = checkToken(tknStr, []byte(jwtKeystring))
	if err != nil {
		log.Println("checkGuiaccess: token error", err.Error())
		return "", err
	}

	return "ok", nil
}

func checkCookie(r *http.Request, name string) (string, error) {
	// Read the cookie as normal.
	c, err := r.Cookie(name)
	if err != nil {
		if err == http.ErrNoCookie {
			log.Println("checkCookie: no cookie:", err.Error())
			return "", errors.New("No cookie found.")
		}
		// For any other type of error, return a bad request status
		log.Println("checkCookie: error:", err.Error())
		return "", err
	}

	// Get the JWT string from the cookie
	log.Println("checkCookie: Value found of", c.Value)
	return c.Value, nil
}

func generateToken(username string, jwtKeystring string, minutes time.Duration) (string, time.Time, error) {
	log.Print("generateToken:", username, "for", minutes, "minutes")
	// get an expiration time of minutes minutes from now
	expirationTime := time.Now().Add(minutes * time.Minute)

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
		log.Println("generateToken: Error creating tokenstring -->", err.Error())
		return "fail", expirationTime, err
	}

	return tokenString, expirationTime, nil
}

func generateApitoken(username string, jwtKeystring string, days int) (string, error) {
	log.Print("generateApitoken:", username, "for", days, "days")

	// get an expiration time days days from now
	expirationTime := time.Now().AddDate(0, 0, days)

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
		log.Println("generateApitoken: Error creating tokenstring -->", err.Error())
		return "fail", err
	}

	return tokenString, nil
}

func checkToken(tokenstr string, jwtKey []byte) (string, error) {
	log.Println("checkToken: checking token", tokenstr)

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
			log.Println("checkToken: SignatureInvalid:", err.Error())
			return "fail", err
		}

		log.Println("checkToken: Token error:", err.Error())
		return "fail", err
	}

	if !tkn.Valid {
		log.Println("checkToken: Token not valid:", tkn.Valid)
		return "fail", errors.New("Token is not valid")
	}

	log.Println("checkToken: valid token exists for user:", claims.Username)
	return claims.Username, nil
}
