package _view

import (
	"cfasuite/internal/_components"
	"cfasuite/internal/_middleware"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"
)


type ViewBuilder struct {
	Title string
	ViewComponents []string
}

func NewViewBuilder(title string, viewComponents []string) *ViewBuilder {
	return &ViewBuilder{
		Title: title,
		ViewComponents: viewComponents,
	}
}

func (v *ViewBuilder) Build() []byte {
	viewString := strings.Join(v.ViewComponents, "");
	return []byte(fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<script src="https://unpkg.com/htmx.org@1.9.10"></script>
			<script src="//unpkg.com/alpinejs" defer></script>
			<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=0">
			<link rel="stylesheet" href="/static/css/output.css">
			<style>
				[x-cloak] { display: none !important; }
			</style>
			<title>%s</title>
		</head>
		<body class='bg-black text-white'>
			%s
		</body>
	`, v.Title, viewString))
}


func Home(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		ctx := map[string]interface{}{}
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		_middleware.MiddlewareChain(ctx, w, r, _middleware.Init,
			func(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
				b := NewViewBuilder("Repairs Log - Login", []string{
					_components.Banner(false, ""),
					_components.MqWrapperCentered(_components.LoginForm(r.URL.Query().Get("err"))),

				})
				w.Write(b.Build())
			},
			_middleware.Log,
		)
	})
}

func Login(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		ctx := map[string]interface{}{}
		_middleware.MiddlewareChain(ctx, w, r, _middleware.Init, _middleware.ParseForm,
			func(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
				username := r.Form.Get("username")
				password := r.Form.Get("password")
				if username == os.Getenv("APP_USERNAME") && password == os.Getenv("APP_PASSWORD") {
					http.SetCookie(w, &http.Cookie{
						Name: "SessionToken",
						Value: os.Getenv("ADMIN_SESSION_TOKEN"),
						Expires: time.Now().Add(24 * time.Hour),
						HttpOnly: true,
					})
					http.Redirect(w, r, "/app", http.StatusSeeOther)
					return
				}
				http.Redirect(w, r, "/?err=invalid credentials", http.StatusSeeOther)
			},
			_middleware.Log,
		)
	})
}

func App(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app", func(w http.ResponseWriter, r *http.Request) {
		ctx := map[string]interface{}{}
		_middleware.MiddlewareChain(ctx, w, r, _middleware.Init, _middleware.Auth,
			func(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
				b := NewViewBuilder("Repairs Log - App", []string{
					_components.Banner(true, ""),
				})
				w.Write(b.Build())
			},
			_middleware.Log,
		)
	})
}




