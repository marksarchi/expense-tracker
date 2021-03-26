package auth

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	//"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// Key is used to store/retrieve a Claims value from a context.Context.
const Key ctxKey = 1

type Keys map[string]*rsa.PrivateKey

// Claims represents the authorization claims transmitted via a JWT.
type Claims struct {
	jwt.StandardClaims
}

type Auth struct {
	mu        sync.RWMutex
	algorithm string
	method    jwt.SigningMethod
	parser    *jwt.Parser
	keys      Keys
}

func New(algorithm string) (*Auth, error) {
	method := jwt.GetSigningMethod(algorithm)
	if method == nil {
		return nil, errors.Errorf("unknown algorithm %v", algorithm)
	}

	// Create the token parser to use. The algorithm used to sign the JWT must be
	// validated to avoid a critical vulnerability:
	// https://auth0.com/blog/critical-vulnerabilities-in-json-web-token-libraries/
	parser := jwt.Parser{
		ValidMethods: []string{algorithm},
	}

	a := Auth{
		algorithm: algorithm,
		method:    method,
		parser:    &parser,
	}

	return &a, nil
}

// GenerateToken generates a signed JWT token string representing the user Claims.
func (a *Auth) GenerateToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(a.method, claims)
	var privateKey = "MIIEpQIBAAKCAQEAvMAHb0IoLvoYuW2kA+LTmnk+hfnBq1eYIh4CT/rMPCxgtzjq"
	// token.Header["kid"] = kid

	// var privateKey *rsa.PrivateKey
	// a.mu.RLock()
	// {
	// 	var ok bool
	// 	privateKey, ok = a.keys[kid]
	// 	if !ok {
	// 		return "", errors.New("kid lookup failed")
	// 	}
	// }
	// a.mu.RUnlock()

	str, err := token.SignedString(privateKey)
	if err != nil {
		return "", errors.Wrap(err, "signing token")
	}

	return str, nil
}

// ValidateToken recreates the Claims that were used to generate a token. It
// verifies that the token was signed using our key.
// func (a *Auth) ValidateToken(tokenStr string) (Claims, error) {

// 	var claims Claims
// 	token, err := a.parser.ParseWithClaims(tokenStr, &claims, a.keyFunc)
// 	if err != nil {
// 		return Claims{}, errors.Wrap(err, "parsing token")
// 	}

// 	if !token.Valid {
// 		return Claims{}, errors.New("invalid token")
// 	}

// 	return claims, nil
// }

func CreateToken(claims jwt.Claims) (string, error) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	// atClaims := jwt.MapClaims{}
	// atClaims["authorized"] = true
	// atClaims["user_id"] = userId
	// atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func VerifyToken(r *http.Request) (jwt.Claims, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims, nil
}
func (a *Auth) ValidateToken(tokenStr string) (jwt.Claims, error) {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims, nil
}
