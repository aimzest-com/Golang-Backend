package handlers

import (
    "net/http"
)

func Main(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte( "<div> <a href='/auth/facebook/login'>Login with Facebook</a> </div> ")) 
}
