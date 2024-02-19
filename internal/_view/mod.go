package _view

import (
	"cfasuite/internal/_components"
	"cfasuite/internal/_middleware"
	"cfasuite/internal/_model"
	"cfasuite/internal/_util"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

type ViewBuilder struct {
	Title          string
	ViewComponents []string
}

func NewViewBuilder(title string, viewComponents []string) *ViewBuilder {
	return &ViewBuilder{
		Title:          title,
		ViewComponents: viewComponents,
	}
}

func (v *ViewBuilder) Build() []byte {
	viewString := strings.Join(v.ViewComponents, "")
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
			<link rel="stylesheet" href="/static/css/animate.css">
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

func (v *ViewBuilder) BuildComponent() []byte {
	viewString := strings.Join(v.ViewComponents, "")
	return []byte(fmt.Sprintf(`%s`, viewString))
}

func ServeStaticFilesAndFavicon(mux *http.ServeMux) {
	mux.HandleFunc("GET /static/", ServeStaticFiles)
	mux.HandleFunc("GET /favicon.ico", ServeFavicon)
}

func ServeFavicon(w http.ResponseWriter, r *http.Request) {
	filePath := "favicon.ico"
	fullPath := filepath.Join(".", "static", filePath)
	http.ServeFile(w, r, fullPath)
}

func ServeStaticFiles(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[len("/static/"):]
	fullPath := filepath.Join(".", "static", filePath)
	http.ServeFile(w, r, fullPath)
}

func PublicTickets(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/ticket/public", func(w http.ResponseWriter, r *http.Request) {
		ctx := map[string]interface{}{}
		_middleware.MiddlewareChain(ctx, w, r, _middleware.Init,
			func(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
				token := r.URL.Query().Get("publicSecurityToken")
				if token != os.Getenv("PUBLIC_SECURITY_TOKEN") {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
				b := NewViewBuilder("Repairs Log - Public Tickets", []string{
					_components.Banner(false, ""),
					_components.Root(
						_components.CenterContentWrapper(
							_components.CreateTicketForm(r.URL.Query().Get("err")),
						),
					),
					_components.Footer(),
				})
				w.Write(b.Build())
			},
			_middleware.Log,
		)
	})

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
					_components.Root(
						_components.CenterContentWrapper(
							_components.LoginForm(r.URL.Query().Get("err"), r.URL.Query().Get("username"), r.URL.Query().Get("password")),
						),
					),
					_components.Footer(),
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
				fmt.Println(username, password)
				fmt.Println(os.Getenv("ADMIN_USERNAME"), os.Getenv("ADMIN_PASSWORD"))
				if username == os.Getenv("ADMIN_USERNAME") && password == os.Getenv("ADMIN_PASSWORD") {
					http.SetCookie(w, &http.Cookie{
						Name:     "SessionToken",
						Value:    os.Getenv("ADMIN_SESSION_TOKEN"),
						Expires:  time.Now().Add(24 * time.Hour),
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
					_components.Root(
						_components.CenterContentWrapper(
							_components.CreateManufacturerForm(r.URL.Query().Get("err"), r.URL.Query().Get("name"), r.URL.Query().Get("phone"), r.URL.Query().Get("email"), ""),
							_components.ManufacturerList(manufacturers, ""),
						),
					),
					_components.Footer(),
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
					Name:  name,
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
					Name:     "SessionToken",
					Value:    "",
					Expires:  time.Now(),
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
				db.Preload("Equipment").First(&manufacturer, id)
				b := NewViewBuilder("Repairs Log - App", []string{
					_components.Banner(true, _components.AppNavMenu(r.URL.Path)),
					_components.Root(
						_components.CenterContentWrapper(
							_components.ManufacturerDetails(manufacturer, r.URL.Query().Get("err")),
							_components.CreateEquipmentForm(r.URL.Query().Get("equipmentErr"), ""),
							_components.EquipmentList(manufacturer.Equipment, ""),
						),
					),
					_components.Footer(),
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
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer/"+id, "err", "invalid delete name provided"), http.StatusSeeOther)
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
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer/"+id, "err", "all fields required"), http.StatusSeeOther)
					return
				}
				if _util.IsValidPhoneNumber(phone) == false {
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer/"+id, "err", "invalid phone format"), http.StatusSeeOther)
					return
				}
				if _util.IsValidEmail(email) == false {
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer/"+id, "err", "invalid email format"), http.StatusSeeOther)
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

func CreateEquipment(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /app/manufacturer/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx := map[string]interface{}{}
		_middleware.MiddlewareChain(ctx, w, r, _middleware.Init, _middleware.ParseMultipartForm, _middleware.Auth,
			func(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				name := r.Form.Get("nickname")
				number := r.Form.Get("number")
				photo, _, err := r.FormFile("photo")
				if err != nil {
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer/"+id, "equipmentErr", "photo required"), http.StatusSeeOther)
					return
				}
				defer photo.Close()
				if photo == nil {
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer/"+id, "equipmentErr", "photo required"), http.StatusSeeOther)
					return
				}
				photoBytes, err := io.ReadAll(photo)
				if err != nil {
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer/"+id, "equipmentErr", "photo required"), http.StatusSeeOther)
					return
				}
				var manufacturer _model.Manufacturer
				db.First(&manufacturer, id)
				qrCodeToken := _util.GenerateRandomToken(32)
				equipment := _model.Equipment{
					Nickname:       name,
					SerialNumber:   number,
					Photo:          photoBytes,
					ManufacturerID: manufacturer.ID,
					QRCodeToken:    qrCodeToken,
				}
				db.Create(&equipment)
				http.Redirect(w, r, "/app/manufacturer/"+id, http.StatusSeeOther)
			},
			_middleware.Log,
		)
	})
}

func Equipment(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/equipment/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx := map[string]interface{}{}
		_middleware.MiddlewareChain(ctx, w, r, _middleware.Init, _middleware.Auth,
			func(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				var equipment _model.Equipment
				var manufacturer _model.Manufacturer
				db.First(&equipment, id)
				db.First(&manufacturer, equipment.ManufacturerID)
				b := NewViewBuilder("Repairs Log - App", []string{
					_components.Banner(true, _components.AppNavMenu(r.URL.Path)),
					_components.Root(
						_components.CenterContentWrapper(
							_components.EquipmentDetails(equipment, manufacturer, r.URL.Query().Get("err")),
							_components.EquipmentQrCodeDownload(equipment),
						),
					),
					_components.Footer(),
				})
				w.Write(b.Build())
			},
			_middleware.Log,
		)
	})
}

func UpdateEquipment(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /app/equipment/{id}/update", func(w http.ResponseWriter, r *http.Request) {
		ctx := map[string]interface{}{}
		_middleware.MiddlewareChain(ctx, w, r, _middleware.Init, _middleware.ParseMultipartForm, _middleware.Auth,
			func(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				name := r.Form.Get("nickname")
				number := r.Form.Get("number")

				// Check if a new photo is provided
				photo, _, err := r.FormFile("photo")
				defer func() {
					if photo != nil {
						photo.Close()
					}
				}()
				if err != nil && err != http.ErrMissingFile {
					http.Redirect(w, r, _util.URLBuilder("/app/equipment/"+id, "updateErr", "failed to retrieve file"), http.StatusSeeOther)
					return
				}

				var equipment _model.Equipment
				db.First(&equipment, id)

				// Update only if a new photo is provided
				if photo != nil {
					photoBytes, err := io.ReadAll(photo)
					if err != nil {
						http.Error(w, "Failed to read file", http.StatusBadRequest)
						return
					}
					equipment.Photo = photoBytes
				}

				equipment.Nickname = name
				equipment.SerialNumber = number
				db.Save(&equipment)
				http.Redirect(w, r, "/app/equipment/"+id, http.StatusSeeOther)
			},
			_middleware.Log,
		)
	})
}

func DeleteEquipment(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /app/equipment/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		ctx := map[string]interface{}{}
		_middleware.MiddlewareChain(ctx, w, r, _middleware.Init, _middleware.ParseForm, _middleware.Auth,
			func(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				name := strings.ToLower(r.Form.Get("name"))
				var equipment _model.Equipment
				db.First(&equipment, id)
				var manufacturer _model.Manufacturer
				db.First(&manufacturer, equipment.ManufacturerID)
				if name != strings.ToLower(equipment.Nickname) {
					http.Redirect(w, r, _util.URLBuilder(fmt.Sprintf("/app/equipment/%d", equipment.ID), "err", "invalid delete name provided"), http.StatusSeeOther)
					return
				}
				db.Delete(&equipment, id)
				http.Redirect(w, r, fmt.Sprintf("/app/manufacturer/%d", manufacturer.ID), http.StatusSeeOther)
			},
			_middleware.Log,
		)
	})
}

func GetEquipmentQRCode(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/equipment/{id}/downloadticketqr", func(w http.ResponseWriter, r *http.Request) {
		ctx := map[string]interface{}{}
		_middleware.MiddlewareChain(ctx, w, r, _middleware.Init, _middleware.Auth, _middleware.IncludePNG,
			func(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				var equipment _model.Equipment
				db.First(&equipment, id)
				url := fmt.Sprintf("%s/app/equipment/%d/ticketform?equipmentToken=%s", os.Getenv("BASE_URL"), equipment.ID, equipment.QRCodeToken)
				// Generate QR code directly into the response writer
				png, err := qrcode.Encode(url, qrcode.Medium, 256)
				if err != nil {
					http.Redirect(w, r, _util.URLBuilder(fmt.Sprintf("/app/equipment/%d", equipment.ID), "err", "failed to generate QR code"), http.StatusSeeOther)
					return
				}
				w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.png"`, equipment.Nickname))
				_, err = w.Write(png)
				if err != nil {
					http.Redirect(w, r, _util.URLBuilder(fmt.Sprintf("/app/equipment/%d", equipment.ID), "err", "failed to write QR code"), http.StatusSeeOther)
					return
				}
			},
			_middleware.Log,
		)
	})
}

func EquipmentTicket(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/equipment/{id}/ticketform", func(w http.ResponseWriter, r *http.Request) {
		ctx := map[string]interface{}{}
		_middleware.MiddlewareChain(ctx, w, r, _middleware.Init,
			func(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				equipmentToken := r.URL.Query().Get("equipmentToken")
				var equipment _model.Equipment
				if err := db.Where("qr_code_token = ?", equipmentToken).First(&equipment).Error; err != nil {
					// Equipment with the provided token does not exist
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
				if strconv.Itoa(int(equipment.ID)) != id {
					// Equipment ID in the URL does not match the retrieved equipment ID
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
				var manufacturer _model.Manufacturer
				db.First(&manufacturer, equipment.ManufacturerID)
				b := NewViewBuilder("Repairs Log - App", []string{
					_components.Banner(true, _components.AppNavMenu(r.URL.Path)),
					_components.Root(
						_components.CenterContentWrapper(
							_components.EquipmentDetails(equipment, manufacturer, r.URL.Query().Get("err")),
							_components.EquipmentQrCodeDownload(equipment),
						),
					),
					_components.Footer(),
				})
				w.Write(b.Build())
			},
			_middleware.Log,
		)
	})
}
