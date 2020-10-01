package main

import (
    "log"
    "net/http"
     "github.com/go-redis/redis/v7"
//     "github.com/twinj/uuid"

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
        "/oauth2/facebook/login",
        handlers.FacebookOauth2Login,
    },
    app.Route {
        "AuthFacebookCallback",
        []string{"GET"},
        "/oauth2/facebook/callback",
        handlers.FacebookOauth2Callback,
    },
    app.Route {
        "AuthRegister",
        []string{"POST"},
        "/auth/register",
        handlers.AuthRegister,
    },
    app.Route {
        "AuthLogin",
        []string{"POST"},
        "/auth/login",
        handlers.AuthLogin,
    },
}

func main() {
    config, err := app.NewConfig()
    if err != nil {
        log.Fatal(err)
    }

    db, err := gorm.Open(sqlite.Open(config.GetString("db_name")), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    redisClient := redis.NewClient(&redis.Options{
        Addr: config.GetString("redis_dsn"),
    })
    _, err = redisClient.Ping().Result()
    if err != nil {
        log.Fatal(err)
    }

    myApp := app.NewApp(config, db)
    router := myApp.NewRouter(routes)

    port := config.GetString("port")
    http.ListenAndServe(":" + port, router)
}
