package _api

import (
	"cfasuite/internal/_middleware"
	"cfasuite/internal/_model"
	"cfasuite/internal/_util"
	"fmt"
	"io"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

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

func CreateTicket(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /app/ticket/public", func(w http.ResponseWriter, r *http.Request) {
		ctx := map[string]interface{}{}
		_middleware.MiddlewareChain(ctx, w, r, _middleware.Init, _middleware.ParseMultipartForm,
			func(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
				securityToken := r.Form.Get("publicSecurityToken")
				creator := r.Form.Get("creator")
				item := r.Form.Get("item")
				problem := r.Form.Get("problem")
				location := r.Form.Get("location")
				if creator == "" || item == "" || problem == "" || location == "" {
					http.Redirect(w, r, _util.URLBuilder("/app/ticket/public", "publicSecurityToken", securityToken, "err", "all fields required", "creator", creator, "item", item, "problem", problem, "location", location), http.StatusSeeOther)
					return
				}
				photo, _, err := r.FormFile("photo")
				if err != nil {
					http.Redirect(w, r, _util.URLBuilder("/app/ticket/public", "publicSecurityToken", securityToken, "err", "photo required", "creator", creator, "item", item, "problem", problem, "location", location), http.StatusSeeOther)
					return
				}
				defer photo.Close()
				if photo == nil {
					http.Redirect(w, r, _util.URLBuilder("/app/ticket/public", "publicSecurityToken", securityToken, "err", "photo required", "creator", creator, "item", item, "problem", problem, "location", location), http.StatusSeeOther)
					return
				}
				photoBytes, err := io.ReadAll(photo)
				if err != nil {
					http.Redirect(w, r, _util.URLBuilder("/app/ticket/public", "publicSecurityToken", securityToken, "err", "photo required", "creator", creator, "item", item, "problem", problem, "location", location), http.StatusSeeOther)
					return
				}
				ticket := _model.Ticket{
					Creator:     creator,
					Item:        item,
					Problem:     problem,
					Location:    location,
					Photo:       photoBytes,
				}
				db.Create(&ticket)
				http.Redirect(w, r, _util.URLBuilder("/app/ticket/public", "publicSecurityToken", securityToken, "success", "your ticket has been created, thank you!"), http.StatusSeeOther)
			},
			_middleware.Log,
		)
	})
}