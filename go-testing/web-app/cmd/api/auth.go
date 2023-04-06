package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const jwtTokenExpiry = time.Minute * 25
const refreshTokenExpiry = time.Hour * 20

type TokenPairs struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// additional settings for tokens
type Claims struct {
	UserName string `json:"name"`
	jwt.RegisteredClaims
}

/*
*	read authorization header
*	extract token
*	verify token
*	return token, claims, error if any
**/
func (app *application) getTokenFromHeaderAndVerify(w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
	// add a header
	w.Header().Add("Vary", "Authorization")

	//we expect auth header to look like: Bearer <token>

	// get the authorization header
	authHeader := r.Header.Get("Authorization")

	// sanity check
	if authHeader == "" {
		return "", nil, errors.New("No auth header.")
	}

	// for expected format: split header on spaces
	headerParts := strings.Split(authHeader, " ")
	// expects only Bearer and the Token
	if len(headerParts) != 2 {
		return "", nil, errors.New("Invalid auth header.")
	}

	// check if we have the word "Bearer"
	if headerParts[0] != "Bearer" {
		return "", nil, errors.New("Unauthorized: no Bearer.")
	}

	token := headerParts[1]

	// declare empty Claims variable
	claims := &Claims{}

	/*	parse token grabbed from auth header
	*	with our claims (read into claims),
	*	using our secret (from the receiver)
	**/
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		//validate the signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(app.JWTSecret), nil
	})

	// check for an error; this also catches expired tokens
	if err != nil {
		if strings.HasPrefix(err.Error(), "Token is expired by") {
			return "", nil, errors.New("Expired token")
		}
		return "", nil, err
	}

	// assure that WE issued this token
	if claims.Issuer != app.Domain {
		return "", nil, errors.New("incorrect token issuer")
	}

	// valid token
	return token, claims, nil
}
