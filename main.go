package main

import (
    "log"
    "net/http"
    "github.com/spf13/viper"
    "github.com/gorilla/mux"

    "backend/app/handlers"
)

func main() {
    loadConfig()

    router := mux.NewRouter()

    router.HandleFunc("/", handlers.Main)
    router.HandleFunc("/auth/facebook/login", handlers.FacebookOauth2Login)
    router.HandleFunc("/auth/facebook/callback", handlers.FacebookOauth2Callback)

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
