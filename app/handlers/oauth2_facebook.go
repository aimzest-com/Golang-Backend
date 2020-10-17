package handlers

import (
    "net/http"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/facebook"

    "backend/app"
    "backend/app/helpers"
)

func FacebookOauth2Login(appContext *app.AppContext, w http.ResponseWriter, r *http.Request) {
    config := appContext.Config
    redirectUrl := config.GetString("host") + ":" + config.GetString("port") + "/auth/facebook/callback"

    facebookOauthConfig := &oauth2.Config{
        RedirectURL: redirectUrl,
        ClientID: config.GetString("facebookClientID"),
        ClientSecret: config.GetString("facebookClientSecret"), 
        Endpoint:     facebook.Endpoint,
    }

    oauthState := helpers.GenerateCookie() 
    session, err := appContext.SessionStore.Get(r, oauthState)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    session.Values["oauth_state"] = oauthState
    err = session.Save(r, w)
    if err != nil {
        http.Error(w, err.Error(), 500)
    }

    fbLoginUrl := facebookOauthConfig.AuthCodeURL(oauthState)
    http.Redirect(w, r, fbLoginUrl, http.StatusTemporaryRedirect)
}

func FacebookOauth2Callback(appContext *app.AppContext, w http.ResponseWriter, r *http.Request) {
    facebookStates, ok :=  r.URL.Query()["state"]
    if !ok {
        http.Error(w, "state query param is not setted", 400)
        return
    }
    facebookState := facebookStates[0]

    session, err := appContext.SessionStore.Get(r, facebookState)
    if err != nil {
        http.Error(w, err.Error(), 400)
        return
    }

    if session.Values["oauth_state"] != facebookState {
        http.Error(w, "oauth state on server is diffent of that delivered by facebook", 400)
    }

    //todo get facebook info about user
    //todo what should we do with the user info. How should it be integrated with other third party authentication providers

    //todo Create user. What's next?
    w.Write([]byte("facebook callback"))
}
