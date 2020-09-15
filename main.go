package main

import (
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/facebook"
    "github.com/gorilla/mux"
    "net/http"
    "fmt"
)

var facebookOauthConfig = &oauth2.Config{
    RedirectURL: "http://localhost:7777/auth/facebook/callback", //todo move url in config
    ClientID: "1010600019384206", //todo move clientId in config
    ClientSecret: "ab3cffd0e090e8b3956357fa0f07320b", //todo move clientsecret in config
    Endpoint:     facebook.Endpoint,
}

func main() {
    router := mux.NewRouter()

    router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte( "<div> <a href='/auth/facebook/login'>Login with Facebook</a> </div> ")) //todo do we need this html?
    })

    router.HandleFunc("/auth/facebook/login", func(w http.ResponseWriter, r *http.Request) {
        oauthState := "oauthState" //todo save oauthState. read docs
        fbLoginUrl := facebookOauthConfig.AuthCodeURL(oauthState)
        fmt.Printf("%s\n", fbLoginUrl)
        http.Redirect(w, r, fbLoginUrl, http.StatusTemporaryRedirect)
    })

    router.HandleFunc("/auth/facebook/callback", func(w http.ResponseWriter, r *http.Request) {
        //todo check for oauth state
        //todo get facebook info about user
        //todo what should we do with the user info. How should it be integrated with other third party authentication providers

        //todo Create user. What's next?
        w.Write([]byte("facebook callback"))
    })


    http.ListenAndServe(":7777", router)//todo move port in a config
}
