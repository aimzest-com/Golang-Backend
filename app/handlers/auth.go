package handlers

import (
    "net/http"
    "fmt"

    "golang.org/x/oauth2"
    "github.com/spf13/viper"
    "golang.org/x/oauth2/facebook"
)

func FacebookOauth2Login(w http.ResponseWriter, r *http.Request) {
    port := viper.GetString("port")
    host := viper.GetString("host")

    facebookOauthConfig := &oauth2.Config{
        RedirectURL: host + ":" + port + "/auth/facebook/callback",
        ClientID: viper.GetString("facebookClientID"), 
        ClientSecret: viper.GetString("facebookClientSecret"), 
        Endpoint:     facebook.Endpoint,
    }

    oauthState := "oauthState" //todo save oauthState. read docs
    fbLoginUrl := facebookOauthConfig.AuthCodeURL(oauthState)
    fmt.Printf("%s\n", fbLoginUrl)
    http.Redirect(w, r, fbLoginUrl, http.StatusTemporaryRedirect)
}

func FacebookOauth2Callback(w http.ResponseWriter, r *http.Request) {
        //todo check for oauth state
        //todo get facebook info about user
        //todo what should we do with the user info. How should it be integrated with other third party authentication providers

        //todo Create user. What's next?
        w.Write([]byte("facebook callback"))
    }
