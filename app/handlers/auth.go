package handlers

import (
    "net/http"
    "golang.org/x/oauth2"
    "github.com/spf13/viper"
    "golang.org/x/oauth2/facebook"

    "crypto/rand"
    "encoding/base64"

    "backend/app"
)

func FacebookOauth2Login(appContext *app.Context, w http.ResponseWriter, r *http.Request) {
    port := viper.GetString("port")
    host := viper.GetString("host")

    facebookOauthConfig := &oauth2.Config{
        RedirectURL: host + ":" + port + "/auth/facebook/callback",
        ClientID: viper.GetString("facebookClientID"), 
        ClientSecret: viper.GetString("facebookClientSecret"), 
        Endpoint:     facebook.Endpoint,
    }

    oauthState := generateStateOauthCookie() //todo 1. store it in a session
    //todo 2. look for a context to be shared through requests
    fbLoginUrl := facebookOauthConfig.AuthCodeURL(oauthState)
    http.Redirect(w, r, fbLoginUrl, http.StatusTemporaryRedirect)
}

func FacebookOauth2Callback(appContext *app.Context, w http.ResponseWriter, r *http.Request) {
        //todo check for oauth state
        //todo get facebook info about user
        //todo what should we do with the user info. How should it be integrated with other third party authentication providers

        //todo Create user. What's next?
        w.Write([]byte("facebook callback"))
}

func generateStateOauthCookie() string {
    b := make([]byte, 128)
    rand.Read(b)
    state := base64.URLEncoding.EncodeToString(b)
    return state
}
