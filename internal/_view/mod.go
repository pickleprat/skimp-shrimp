package _view

import (
	"cfasuite/internal/_components"
	"cfasuite/internal/_middleware"
	"fmt"
	"net/http"
	"strings"

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
			<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=0">
			<link rel="stylesheet" href="/static/css/output.css">
			<title>%s</title>
		</head>
		<body class='bg-black text-white'>
			%s
		</body>
	`, v.Title, viewString))
}


func Home(db *gorm.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := map[string]interface{}{}
		if r.Method == "GET" {
			_middleware.MiddlewareChain(ctx, w, r, _middleware.Init,
				func(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
					b := NewViewBuilder("Repairs Log - Login", []string{
						_components.Banner(false),
						_components.MqWrapperCentered(_components.LoginForm(r.URL.Query().Get("err"))),
	
					})
					view := b.Build()
					w.Write(view)
				},
				_middleware.Log,
			)
		}
	})
}




