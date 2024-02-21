package _middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func MiddlewareChain(
	ctx map[string]interface{}, 
	w http.ResponseWriter, 
	r *http.Request, 
	handler func(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request), 
	middleware ...func(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) error) {
	for _, mw := range middleware {
		err := mw(ctx, w, r)
		if err != nil {
			return
		}
	}
	handler(ctx, w, r)
	Log(ctx, w, r)
}


func Log(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) error {
	startTime := ctx["StartTime"].(time.Time)
	elapsedTime := time.Since(startTime)
	formattedTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] [%s] [%s] [%s]\n", formattedTime, r.Method, r.URL.Path, elapsedTime)
	return nil
}

func Init(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) error {
	ctx["StartTime"] = time.Now()
	return nil
}

func SvgHeaders(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "image/svg+xml")
	return nil
}


func ParseForm(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()
	return nil
}

func ParseMultipartForm(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) error {
	r.ParseMultipartForm(10 << 20)
	return nil
}

func Auth(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) error {
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

func IncludePNG(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "image/png")
	return nil
}
