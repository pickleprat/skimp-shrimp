package _view

import (
	"cfasuite/internal/_components"
	"cfasuite/internal/_middleware"
	"cfasuite/internal/_model"
	"cfasuite/internal/_util"
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
			<script src="/static/js/index.js" defer></script>
			<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=0">
			<link rel="stylesheet" href="/static/css/output.css">
			<style>
				[x-cloak] { display: none !important; }
			</style>
			<title>%s</title>
		</head>
		<body hx-boost='true' class='bg-black text-white'>
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
					_components.MqWrapperCentered(_components.LoginForm(r.URL.Query().Get("err"), r.URL.Query().Get("username"), r.URL.Query().Get("password"))),

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
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
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
				http.Redirect(w, r, _util.URLBuilder("/", "err", "invalid credentials", "username", username), http.StatusSeeOther)
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
				var manufacturers []_model.Manufacturer
                db.Find(&manufacturers)
				b := NewViewBuilder("Repairs Log - App", []string{
					_components.Banner(true, _components.AppNavMenu(r.URL.Path)),
					_components.MqGridTwoColEvenSplit(
						_components.CreateManufacturerForm(r.URL.Query().Get("err"), r.URL.Query().Get("name"), r.URL.Query().Get("phone"), r.URL.Query().Get("email"), ""),
						_components.ManufacturerList(manufacturers, ""),
					),
				})
				w.Write(b.Build())
			},
			_middleware.Log,
		)
	})
}

func CreateManufacturer(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /app", func(w http.ResponseWriter, r *http.Request) {
		ctx := map[string]interface{}{}
		_middleware.MiddlewareChain(ctx, w, r, _middleware.Init, _middleware.Auth, _middleware.ParseForm,
			func(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
				name := r.Form.Get("name")
				phone := r.Form.Get("phone")
				email := r.Form.Get("email")
				if name == "" || phone == "" || email == "" {
					http.Redirect(w, r, _util.URLBuilder("/app", "err", "all fields required", "name", name, "phone", phone, "email", email), http.StatusSeeOther)
					return
				}
				if _util.IsValidPhoneNumber(phone) == false {
					http.Redirect(w, r, _util.URLBuilder("/app", "err", "invalid phone format", "name", name, "phone", phone, "email", email), http.StatusSeeOther)
					return
				}
				if _util.IsValidEmail(email) == false {
					http.Redirect(w, r, _util.URLBuilder("/app", "err", "invalid email format", "name", name, "phone", phone, "email", email), http.StatusSeeOther)
					return
				}
				manufacturer := _model.Manufacturer{
					Name: name,
					Phone: phone,
					Email: email,
				}
				db.Create(&manufacturer)
				http.Redirect(w, r, "/app", http.StatusSeeOther)
			},
			_middleware.Log,
		)
	})

}

func Logout(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /logout", func(w http.ResponseWriter, r *http.Request) {
		ctx := map[string]interface{}{}
		_middleware.MiddlewareChain(ctx, w, r, _middleware.Init, _middleware.Auth,
			func(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
				http.SetCookie(w, &http.Cookie{
					Name: "SessionToken",
					Value: "",
					Expires: time.Now(),
					HttpOnly: true,
				})
				http.Redirect(w, r, "/", http.StatusSeeOther)
			},
			_middleware.Log,
		)
	})
}

func Manufacturer(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/manufacturer/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx := map[string]interface{}{}
		_middleware.MiddlewareChain(ctx, w, r, _middleware.Init, _middleware.Auth,
			func(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				var manufacturer _model.Manufacturer
				db.First(&manufacturer, id)
				b := NewViewBuilder("Repairs Log - App", []string{
					_components.Banner(true, _components.AppNavMenu(r.URL.Path)),
					_components.MqGridTwoColEvenSplit(
						_components.ManufacturerDetails(manufacturer, r.URL.Query().Get("update"), r.URL.Query().Get("deleteErr"), r.URL.Query().Get("updateErr"), ""),
						_components.CreateEquipmentForm(r.URL.Query().Get("equipmentErr"), ""),
					),
				})
				w.Write(b.Build())
			},
			_middleware.Log,
		)
	})
}

func DeleteManufacturer(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /app/manufacturer/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		ctx := map[string]interface{}{}
		_middleware.MiddlewareChain(ctx, w, r, _middleware.Init, _middleware.ParseForm, _middleware.Auth,
			func(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				name := strings.ToLower(r.Form.Get("name"))
				var manufacturer _model.Manufacturer
				db.First(&manufacturer, id)
				if name != strings.ToLower(manufacturer.Name) {
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer/"+id, "deleteErr", "invalid name", "update", "false"), http.StatusSeeOther)
					return
				
				}
				db.Delete(&manufacturer, id)
				http.Redirect(w, r, "/app", http.StatusSeeOther)
			},
			_middleware.Log,
		)
	})
}

func UpdateManufacturer(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /app/manufacturer/{id}/update", func(w http.ResponseWriter, r *http.Request) {
		ctx := map[string]interface{}{}
		_middleware.MiddlewareChain(ctx, w, r, _middleware.Init, _middleware.ParseForm, _middleware.Auth,
			func(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				name := r.Form.Get("name")
				phone := r.Form.Get("phone")
				email := r.Form.Get("email")
				if name == "" || phone == "" || email == "" {
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer/"+id, "updateErr", "all fields required"), http.StatusSeeOther)
					return
				}
				if _util.IsValidPhoneNumber(phone) == false {
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer/"+id, "updateErr", "invalid phone format"), http.StatusSeeOther)
					return
				}
				if _util.IsValidEmail(email) == false {
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer/"+id, "updateErr", "invalid email format"), http.StatusSeeOther)
					return
				}
				var manufacturer _model.Manufacturer
				db.First(&manufacturer, id)
				manufacturer.Name = name
				manufacturer.Phone = phone
				manufacturer.Email = email
				db.Save(&manufacturer)
				http.Redirect(w, r, "/app/manufacturer/"+id, http.StatusSeeOther)
			},
			_middleware.Log,
		)
	})
}




