package handlers

import (
    "net/http"
    "backend/app"
)

func Main(appContext *app.AppContext, w http.ResponseWriter, r *http.Request) {
    w.Write([]byte( "<div> <a href='/auth/facebook/login'>Login with Facebook</a> </div> ")) 
}
