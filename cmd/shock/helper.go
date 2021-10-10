
package main

import (
	"errors"
	"encoding/base64"

	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

)

const (


	authHeaderName      string = "Authorization"
	authHeaderPrefix    string = "Bearer "
	authHeaderPrefixLen int    = len(authHeaderPrefix)
	minLegalTokenLength int    = authHeaderPrefixLen + 5 // a.b.c for a jwt
)

type contextKeyType string

type service struct {
	parser    jwt.Parser
	secret    []byte
	nextPongs map[string]time.Time
	mutex     sync.Mutex
}

func decodeSecret(encodedSecret string) []byte {
	if encodedSecret == "" {
		log.Print("secret is required")
		return []byte{}
	}

	secret, err := base64.StdEncoding.DecodeString(encodedSecret)
	if err != nil {
		log.Printf("Could not parse secret: %v", err)
		return []byte{}
	}

	return secret
}

// newService creates an instance of our service data that stores the secret and JWT parser
func newService(secret []byte) *service {
	return &service{
		parser:    jwt.Parser{ValidMethods: []string{"HS256"}},
		secret:    secret,
		nextPongs: make(map[string]time.Time),
	}
}


func (s *service) getKey(*jwt.Token) (interface{}, error) {
	return s.secret, nil
}

// extract claims from the token
func (s *service) claims(token string) (*jwtClaims, error) {

	parsedToken, err := s.parser.ParseWithClaims(token, &jwtClaims{}, s.getKey)

	if err != nil {
		log.Print(err)
		return &jwtClaims{}, errors.New("Could not parse authorization header")
	}

	if claims, ok := parsedToken.Claims.(*jwtClaims); ok && parsedToken.Valid {
		return claims, nil
	}
	log.Print("Could not parse JWT claims")
	return &jwtClaims{}, errors.New("Could not parse JWT claims")
}

// find token from header otherwise error description
func (s *service) token(r *http.Request) (string, error) {
	var token string

	tokens, ok := r.Header[authHeaderName]
	if !ok {
		log.Print("Missing authorization header")
		return "Missing authorization header", errors.New("missing")
	}

	if len(tokens) != 1 {
		log.Print("Multiple authorization headers found")
		return "Multiple authorization headers found; only one header should be sent", errors.New("multiple")
	}

	token = tokens[0]
	if !strings.HasPrefix(token, authHeaderPrefix) || len(token) < minLegalTokenLength {
		log.Print("Malformed authorization header")
		return "Malformed authorization header", errors.New("malformed")
	}
	token = token[authHeaderPrefixLen:]
	return token, nil
}

