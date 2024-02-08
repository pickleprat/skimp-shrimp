package _view

import (
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
			<title>%s</title>
		</head>
		<body>
			%s
		</body>
	`, v.Title, viewString))
}


func Home(db *gorm.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := map[string]interface{}{}
		_middleware.MiddlewareChain(ctx, w, r, _middleware.Init,
			func(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
				b := NewViewBuilder("CFA Suite - Home", []string{"<h1>CFASuite</h1>"})
				fmt.Println(b.ViewComponents)
				view := b.Build()
				w.Write(view)
			},
			_middleware.Log,
		)
	})
}




