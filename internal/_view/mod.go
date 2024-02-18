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
	viewString := strings.Join(v.ViewComponents, "");
	return []byte(fmt.Sprintf(`%s`, viewString))
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
				fmt.Println(username, password)
				fmt.Println(os.Getenv("APP_USERNAME"), os.Getenv("APP_PASSWORD"))
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
					_components.Root(
						_components.MqGridTwoColEvenSplit(
							_components.CreateManufacturerForm(r.URL.Query().Get("err"), r.URL.Query().Get("name"), r.URL.Query().Get("phone"), r.URL.Query().Get("email"), ""),
							_components.ManufacturerList(manufacturers, ""),
						),
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
				db.Preload("Equipment").First(&manufacturer, id)
				b := NewViewBuilder("Repairs Log - App", []string{
					_components.Banner(true, _components.AppNavMenu(r.URL.Path)),
					_components.Root(
						_components.MqGridTwoColEvenSplit(
							_components.ManufacturerDetails(manufacturer, r.URL.Query().Get("update"), r.URL.Query().Get("deleteErr"), r.URL.Query().Get("updateErr"), ""),
							_components.CreateEquipmentForm(r.URL.Query().Get("equipmentErr"), ""),
						),
						_components.MqGridOneRowCentered(
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
                    http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
                    return
                }
                defer photo.Close()
				photoBytes, err := io.ReadAll(photo)
				if err != nil {
					http.Error(w, "Failed to read file", http.StatusBadRequest)
					return
				}
				var manufacturer _model.Manufacturer
				db.First(&manufacturer, id)
				equipment := _model.Equipment{
					Nickname: name,
					SerialNumber: number,
					Photo: photoBytes,
					ManufacturerID: manufacturer.ID,
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
						_components.MqGridTwoColEvenSplit(
							_components.EquipmentDetails(equipment, manufacturer, r.URL.Query().Get("updateErr")),
							"",
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

func EquipmentSettingsForm(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/equipment/{id}/settings", func(w http.ResponseWriter, r *http.Request) {
		ctx := map[string]interface{}{}
		_middleware.MiddlewareChain(ctx, w, r, _middleware.Init, _middleware.ComponentAuth,
			func(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				var equipment _model.Equipment
				db.First(&equipment, id)
				b := NewViewBuilder("Repairs Log - App", []string{
					fmt.Sprintf(`
						<span id='equipment-settings-form'>
							<div class='bg-black opacity-70 fixed h-full w-full top-0 z-40'></div>					
							<div class='fixed h-full w-full z-50 grid top-0 mt-4 sm:mt-8 px-4'>
								<form enctype='multipart/form-data' method='POST' action='/app/equipment/%d/update' x-data="{ loading: false }" class='fade-in gap-4 md:place-self-center mt-[75px] md:mt-0 mb-[500px] grid w-full p-6 md:max-w-[500px] border border-gray rounded bg-black'>
									<div class='flex justify-between'>
										<h2>Equipment Settings</h2>
										%s
									</div>
									<div class='flex flex-col gap-4'>
										%s%s%s%s%s
									</div>
								</form>
							</div>
						</span>
					`, 
						equipment.ID,
						_components.SvgIcon("/static/svg/x-dark.svg", "sm", "hx-get='/app/component/clear' hx-trigger='click' hx-swap='outerHTML' hx-target='#equipment-settings-form'"),
						_components.FormInputLabel("Nickname", "nickname", "text", equipment.Nickname),
						_components.FormInputLabel("Serial Number", "number", "text", equipment.SerialNumber),
						_components.FormPhotoUpload(),
						_components.FormLoader(),
						_components.FormSubmitButton(),

					),
				})
				w.Write(b.BuildComponent())
			},
			_middleware.Log,
		)
	})
}


func ClearComponent(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/component/clear", func(w http.ResponseWriter, r *http.Request) {
		ctx := map[string]interface{}{}
		_middleware.MiddlewareChain(ctx, w, r, _middleware.Init,
			func(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
				b := NewViewBuilder("Repairs Log - App", []string{``})
				w.Write(b.BuildComponent())
			},
			_middleware.Log,
		)
	})
}

func ClientRedirect(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/component/redirect", func(w http.ResponseWriter, r *http.Request) {
		ctx := map[string]interface{}{}
		_middleware.MiddlewareChain(ctx, w, r, _middleware.Init,
			func(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
				b := NewViewBuilder("Repairs Log - App", []string{`
					<script>
						htmx.ajax("GET", "/", "body");
					</script>
				`})
				w.Write(b.BuildComponent())
			},
			_middleware.Log,
		)
	})
}


