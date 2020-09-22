package app

import (
    "net/http"
    "github.com/gorilla/mux"
    "github.com/gorilla/sessions"
    "gorm.io/gorm"

    "backend/app/model"
)

type Route struct {
    Name string
    Method []string
    Path string
    ContextHandlerFunc ContextHandlerFunc
}

type Routes []Route

type AppContext struct{
    Config *Config
    SessionStore *sessions.CookieStore
    Db *gorm.DB
}

type ContextHandlerFunc func(*AppContext, http.ResponseWriter, *http.Request)

type ContextHandler struct {
    AppContext *AppContext
    ContextHandlerFunc ContextHandlerFunc
}
func (contextHandler ContextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    contextHandler.ContextHandlerFunc(contextHandler.AppContext, w, r)
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
            Handler(&ContextHandler{app.Context, route.ContextHandlerFunc})
    }

    return router
}

func NewApp(config *Config, db *gorm.DB) *App {
    sessionStore := sessions.NewCookieStore([]byte(config.GetString("session_key")))

    db.AutoMigrate(&model.User{})

    return &App{
        Context: &AppContext{
            Config: config,
            SessionStore: sessionStore,
            Db: db,
        },
    }
}
