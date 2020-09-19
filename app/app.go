package app

import (
    "net/http"
)

type AppContext struct{}

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

func (app *App) Bind(ctxHandlerFunc ContextHandlerFunc) http.Handler {
    return &ContextHandler{app.Context, ctxHandlerFunc}
}

func NewApp() *App {
    return &App{}
}
