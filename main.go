package main

import (
    "log"
    "net/http"
    "github.com/spf13/viper"

    "backend/app/handlers"
    "backend/app"
)

var routes = app.Routes{
    app.Route {
        "default",
        []string{"GET"},
        "/",
        handlers.Main,
    },
    app.Route {
        "AuthFacebookLogin",
        []string{"GET"},
        "/auth/facebook/login",
        handlers.FacebookOauth2Login,
    },
    app.Route {
        "AuthFacebookCallback",
        []string{"GET"},
        "/auth/facebook/callback",
        handlers.FacebookOauth2Callback,
    },
}

//todo move config in AppContext
//todo move sessions in AppContext

func main() {
    loadConfig()

    myApp := app.NewApp()
    router := myApp.NewRouter(routes)

    /*
    router.Methods("GET").Path("/").Name("default").Handler(myApp.Bind(handlers.Main))
    router.Methods("GET").Path("/auth/facebook/login").Name("AuthFacebookLogin").Handler(myApp.Bind(handlers.FacebookOauth2Login))
    router.Methods("GET").Path("/auth/facebook/callback").Name("AuthFacebookCallback").Handler(myApp.Bind(handlers.FacebookOauth2Callback))
    */

    port := viper.GetString("port")
    http.ListenAndServe(":" + port, router)
}

func loadConfig() {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    err := viper.ReadInConfig()
    if err != nil {
        log.Fatal(err)
    }
}
