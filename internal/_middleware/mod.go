package _middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func MiddlewareChain(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request, middleware ...func(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request)) {
	for _, mw := range middleware {
		mw(ctx, w, r)
	}
}

func Log(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	startTime := ctx["StartTime"].(time.Time)
	elapsedTime := time.Since(startTime)
	formattedTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] [%s] [%s] [%s]\n", formattedTime, r.Method, r.URL.Path, elapsedTime)
}

func Init(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	ctx["StartTime"] = time.Now()
}

func SvgHeaders(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
}


func ParseForm(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
}

func ParseMultipartForm(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
}

func Auth(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SessionToken")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if cookie.Value != os.Getenv("ADMIN_SESSION_TOKEN") {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func IncludePNG(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
}
