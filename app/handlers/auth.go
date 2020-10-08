package handlers

import (
    "net/http"
    "encoding/json"
    "strings"
    "golang.org/x/crypto/bcrypt"

    "backend/app"
    "backend/app/model"
    "backend/app/form"

    "fmt"
)

func AuthRegister(appContext *app.AppContext, w http.ResponseWriter, r *http.Request) {
    var registerForm form.Register
    err := json.NewDecoder(r.Body).Decode(&registerForm)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err = appContext.Validate.Struct(registerForm)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    username := strings.TrimSpace(registerForm.Username)

    var dbUser model.User
    result := appContext.Db.Where("username = ?", username).Find(&dbUser)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusBadRequest)
        return
    }

    if result.RowsAffected > 0 {
        http.Error(w, "The username is already in use", http.StatusBadRequest)
        return
    }

    password := []byte(strings.TrimSpace(registerForm.Password))
    password, err = bcrypt.GenerateFromPassword(password, appContext.Config.GetInt("bcrypt_cost"))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    user := model.User{Username: username, Password: string(password)}
    result = appContext.Db.Create(&user)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusBadRequest)
        return
    }

    w.Write([]byte(fmt.Sprintf("%d", user.ID)))
}

func AuthLogin(appContext *app.AppContext, w http.ResponseWriter, r *http.Request) {
    var loginForm form.Login
    err := json.NewDecoder(r.Body).Decode(&loginForm)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err = appContext.Validate.Struct(loginForm)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    username := strings.TrimSpace(loginForm.Username)
    var user model.User
    result := appContext.Db.Where("username = ?", username).Find(&user)
    if result.RowsAffected == 0 {
        http.Error(w, "user not found", http.StatusNotFound)
        return
    }

    password := strings.TrimSpace(loginForm.Password)
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        http.Error(w, "wrong password", http.StatusNotFound)
        return
    }

    jwtToken, err := appContext.JWTStorage.NewToken(user.ID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnprocessableEntity)
        return
    }

    err = appContext.JWTStorage.CreateAuth(uint64(user.ID), jwtToken)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    tokens := map[string]string {
        "access_token": jwtToken.AccessToken,
        "refresh_token": jwtToken.RefreshToken,
    }

    tokensJson, err := json.Marshal(tokens)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Write([]byte(tokensJson))
}
