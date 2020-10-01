package helpers

import (
    "math/rand"
    "time"
    "github.com/dgrijalva/jwt-go"
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
var charset string = "!#$%&'*+-.0123456789ABCDEFGHIJKLMNOPQRSTUWVXYZ^_`abcdefghijklmnopqrstuvwxyz|~"
var lenCharset int = len(charset)

func GenerateCookie() string {
    b := make([]byte, 40) //php session id has length 40

    for i := range b {
        b[i] = charset[seededRand.Intn(lenCharset)]
    }

    return string(b)
}

func GenerateJWTToken(userId uint, secret string) (string, error) {
    atClaims := jwt.MapClaims{}
    atClaims["authorized"] = true
    atClaims["user_id"] = userId
    atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()

    at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
    token, err := at.SignedString([]byte(secret))
    if err != nil {
        return "", nil
    }
    return token, nil
}
