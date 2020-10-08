package handlers

import (
    "net/http"
    "backend/app"
)

func Test(appContext *app.AppContext, w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("test handler"))
}
