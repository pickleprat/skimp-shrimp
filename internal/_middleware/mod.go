package _middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type CustomContext struct {
    StartTime time.Time
}

type CustomHandler func(ctx *CustomContext, w http.ResponseWriter, r *http.Request)
type CustomMiddleware func(ctx *CustomContext, w http.ResponseWriter, r *http.Request) error

func MiddlewareChain(w http.ResponseWriter, r *http.Request, handler CustomHandler, middleware ...CustomMiddleware) {
    customContext := &CustomContext{
        StartTime: time.Now(),
    }
    for _, mw := range middleware {
        err := mw(customContext, w, r)
        if err != nil {
            return
        }
    }
    handler(customContext, w, r)
    Log(customContext, w, r)
}

func Log(ctx *CustomContext, w http.ResponseWriter, r *http.Request) error {
    elapsedTime := time.Since(ctx.StartTime)
    formattedTime := time.Now().Format("2006-01-02 15:04:05")
    fmt.Printf("[%s] [%s] [%s] [%s]\n", formattedTime, r.Method, r.URL.Path, elapsedTime)
    return nil
}

func Init(ctx *CustomContext, w http.ResponseWriter, r *http.Request) error {
    ctx.StartTime = time.Now()
    return nil
}

func SvgHeaders(ctx *CustomContext, w http.ResponseWriter, r *http.Request) error {
    w.Header().Set("Content-Type", "image/svg+xml")
    return nil
}

func ParseForm(ctx *CustomContext, w http.ResponseWriter, r *http.Request) error {
    r.ParseForm()
    return nil
}

func ParseMultipartForm(ctx *CustomContext, w http.ResponseWriter, r *http.Request) error {
    r.ParseMultipartForm(10 << 20)
    return nil
}

func Auth(ctx *CustomContext, w http.ResponseWriter, r *http.Request) error {
    cookie, err := r.Cookie("SessionToken")
    if err != nil {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return fmt.Errorf("no session token")
    }
    if cookie.Value != os.Getenv("ADMIN_SESSION_TOKEN") {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return fmt.Errorf("invalid session token")
    }
    return nil
}

func IncludePNG(ctx *CustomContext, w http.ResponseWriter, r *http.Request) error {
    w.Header().Set("Content-Type", "image/png")
    return nil
}
