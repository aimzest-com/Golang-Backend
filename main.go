package main

import (
    "log"
    "net/http"
    "github.com/spf13/viper"
    "github.com/gorilla/mux"

    "backend/app/handlers"
    "backend/app"
)

type ContextHandler struct {
    AppContext *app.Context
    ContextHandlerFunc func(*app.Context, http.ResponseWriter, *http.Request)
}
func (contextHandler ContextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    contextHandler.ContextHandlerFunc(contextHandler.AppContext, w, r)
}

func main() {
    loadConfig()

    appContext := &app.Context{}

    router := mux.NewRouter()

    router.Methods("GET").Path("/").Name("default").Handler(ContextHandler{appContext, handlers.Main})
    router.Methods("GET").Path("/auth/facebook/login").Name("AuthFacebookLogin").Handler(ContextHandler{appContext, handlers.FacebookOauth2Login})
    router.Methods("GET").Path("/auth/facebook/callback").Name("AuthFacebookCallback").Handler(ContextHandler{appContext, handlers.FacebookOauth2Callback})

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
