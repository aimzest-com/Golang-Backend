package main

import (
    "log"
    "net/http"
    "gorm.io/gorm"
    "gorm.io/driver/sqlite"

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

func main() {
    config, err := app.NewConfig()
    if err != nil {
        log.Fatal(err)
    }

    db, err := gorm.Open(sqlite.Open(config.GetString("DB_NAME")), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    myApp := app.NewApp(config, db)
    router := myApp.NewRouter(routes)

    port := config.GetString("port")
    http.ListenAndServe(":" + port, router)
}
