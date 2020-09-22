package handlers

import (
    "net/http"
    "encoding/json"

    "backend/app"
    "backend/app/model"

    "fmt"
)

func AuthRegister(appContext *app.AppContext, w http.ResponseWriter, r *http.Request) {
    var user model.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    appContext.Db.Create(&user)

    //todo add uniqueness for username
    //todo check errors

    w.Write([]byte(fmt.Sprintf("%d", user.ID)))
}
