package app

import (
    "net/http"
    PlaygroundValidate "github.com/go-playground/validator/v10"

    "github.com/gorilla/mux"
    "github.com/gorilla/sessions"
    "gorm.io/gorm"

    "backend/app/model"
    "github.com/go-redis/redis/v7"
)

type RouteOptions struct {
    Authenticate bool
}

type Route struct {
    Name string
    Method []string
    Path string
    ContextHandlerFunc ContextHandlerFunc
    Options *RouteOptions
}

type Routes []Route

type AppContext struct{
    Config *Config
    SessionStore *sessions.CookieStore
    Db *gorm.DB
    Validate *PlaygroundValidate.Validate
    JWTStorage *JWTStorage
}

type ContextHandlerFunc func(*AppContext, http.ResponseWriter, *http.Request)

type ContextHandler struct {
    AppContext *AppContext
    ContextHandlerFunc ContextHandlerFunc
    RouteOptions *RouteOptions
}

func (contextHandler ContextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    contextHandlerFunc := contextHandler.ContextHandlerFunc

    if contextHandler.RouteOptions.Authenticate {
        contextHandlerFunc = authenticateContextHandlerFunc(contextHandlerFunc)
    }

    contextHandlerFunc(contextHandler.AppContext, w, r)
}

func authenticateContextHandlerFunc(contextHandlerFunc ContextHandlerFunc) ContextHandlerFunc {
    return  func(appContext *AppContext, w http.ResponseWriter, r *http.Request) {
        accessDetails, err := appContext.JWTStorage.ExtractTokenMetadata(r)
        if err != nil {
            http.Error(w, err.Error(), http.StatusUnauthorized)
            return;
        }

        _, err = appContext.JWTStorage.FetchAuth(accessDetails)
        if err != nil {
            http.Error(w, err.Error(), http.StatusUnauthorized)
            return;
        }

        contextHandlerFunc(appContext, w, r)
    }
}

type App struct {
    Context *AppContext
}

func (app *App) NewRouter(routes Routes) *mux.Router {
    router := mux.NewRouter()

    for _, route := range routes {
        router.
            Methods(route.Method...).
            Path(route.Path).
            Name(route.Name).
            Handler(&ContextHandler{app.Context, route.ContextHandlerFunc, route.Options})
    }

    return router
}

func NewApp(config *Config, db *gorm.DB, redisClient *redis.Client) *App {
    sessionStore := sessions.NewCookieStore([]byte(config.GetString("session_key")))

    db.AutoMigrate(&model.User{})

    jwtAccessSecret := config.GetString("jwt_access_secret")
    jwtRefreshSecret := config.GetString("jwt_refresh_secret")
    jwtStorage := NewJWTStorage(redisClient, jwtAccessSecret, jwtRefreshSecret)

    return &App{
        Context: &AppContext{
            Config: config,
            SessionStore: sessionStore,
            Db: db,
            Validate: PlaygroundValidate.New(),
            JWTStorage: jwtStorage,
        },
    }
}
