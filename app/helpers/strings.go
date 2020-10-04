package helpers

import (
    "math/rand"
    "time"
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
