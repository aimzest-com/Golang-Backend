package main

import (
    "log"
    "net/http"
    "github.com/spf13/viper"
    "github.com/gorilla/mux"

    "backend/app/handlers"
    "backend/app"
)

func main() {
    loadConfig()

    myApp := app.NewApp()
    router := mux.NewRouter()

    router.Methods("GET").Path("/").Name("default").Handler(myApp.Bind(handlers.Main))
    router.Methods("GET").Path("/auth/facebook/login").Name("AuthFacebookLogin").Handler(myApp.Bind(handlers.FacebookOauth2Login))
    router.Methods("GET").Path("/auth/facebook/callback").Name("AuthFacebookCallback").Handler(myApp.Bind(handlers.FacebookOauth2Callback))

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
