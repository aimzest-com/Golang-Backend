package handlers

import (
    "net/http"
    "encoding/json"
    "strings"
    "golang.org/x/crypto/bcrypt"

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

    err = appContext.Validate.Struct(user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    user.Username = strings.TrimSpace(user.Username)

    var dbUser model.User
    result := appContext.Db.Where("username = ?", user.Username).Find(&dbUser)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusBadRequest)
        return
    }

    if result.RowsAffected > 0 {
        http.Error(w, "The username is already in use", http.StatusBadRequest)
        return
    }

    password := []byte(strings.TrimSpace(user.Password))
    password, err = bcrypt.GenerateFromPassword(password, appContext.Config.GetInt("bcrypt_cost"))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    user.Password = string(password)
    result = appContext.Db.Create(&user)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusBadRequest)
        return
    }

    w.Write([]byte(fmt.Sprintf("%d", user.ID)))
}
