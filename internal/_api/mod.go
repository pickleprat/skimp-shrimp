package _api

import (
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

func Login(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				username := r.Form.Get("username")
				password := r.Form.Get("password")
				if username == os.Getenv("ADMIN_USERNAME") && password == os.Getenv("ADMIN_PASSWORD") {
					http.SetCookie(w, &http.Cookie{
						Name:     "SessionToken",
						Value:    os.Getenv("ADMIN_SESSION_TOKEN"),
						Expires:  time.Now().Add(24 * time.Hour),
						HttpOnly: true,
					})
					http.Redirect(w, r, "/app/ticket", http.StatusSeeOther)
					return
				}
				http.Redirect(w, r, _util.URLBuilder("/", "err", "invalid credentials", "username", username), http.StatusSeeOther)
			},
			_middleware.Init, _middleware.ParseForm,
		)
	})
}

func Logout(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /logout", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				http.SetCookie(w, &http.Cookie{
					Name:     "SessionToken",
					Value:    "",
					Expires:  time.Now(),
					HttpOnly: true,
				})
				http.Redirect(w, r, "/", http.StatusSeeOther)
			},
			_middleware.Init, _middleware.Auth,
		)
	})
}

func CreateManufacturer(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /app/manufacturer", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				name := r.Form.Get("name")
				phone := r.Form.Get("phone")
				email := r.Form.Get("email")
				submitRedirect := r.Form.Get("submitRedirect")
				if name == "" || phone == "" || email == "" {
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer", "err", "all fields required", "name", name, "phone", phone, "email", email), http.StatusSeeOther)
					return
				}
				if _util.IsValidPhoneNumber(phone) == false {
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer", "err", "invalid phone format", "name", name, "phone", phone, "email", email), http.StatusSeeOther)
					return
				}
				if _util.IsValidEmail(email) == false {
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer", "err", "invalid email format", "name", name, "phone", phone, "email", email), http.StatusSeeOther)
					return
				}
				manufacturer := _model.Manufacturer{
					Name:  name,
					Phone: phone,
					Email: email,
				}
				db.Create(&manufacturer)
				redirectURL := _util.ConditionalString(
					submitRedirect == "", 
					_util.URLBuilder("/app/manufacturer", "success", "manufacturer created"), 
					submitRedirect,
				)
				http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			},
			_middleware.Init, _middleware.Auth, _middleware.ParseForm,
		)
	})
}

func DeleteManufacturer(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /app/manufacturer/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				name := strings.ToLower(r.Form.Get("name"))
				var manufacturer _model.Manufacturer
				db.First(&manufacturer, id)
				if name != strings.ToLower(manufacturer.Name) {
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer/"+id, "err", "invalid delete name provided", "form", "delete"), http.StatusSeeOther)
					return

				}
				db.Delete(&manufacturer, id)
				http.Redirect(w, r, _util.URLBuilder("/app/manufacturer", "success", fmt.Sprintf("manufacturer %s deleted", strings.ToLower(manufacturer.Name))), http.StatusSeeOther)
			},
			_middleware.Init, _middleware.ParseForm, _middleware.Auth,
		)
	})
}

func UpdateManufacturer(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /app/manufacturer/{id}/update", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				name := r.Form.Get("name")
				phone := r.Form.Get("phone")
				email := r.Form.Get("email")
				if name == "" || phone == "" || email == "" {
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer/"+id, "err", "all fields required", "form", "update"), http.StatusSeeOther)
					return
				}
				if _util.IsValidPhoneNumber(phone) == false {
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer/"+id, "err", "invalid phone format", "form", "update"), http.StatusSeeOther)
					return
				}
				if _util.IsValidEmail(email) == false {
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer/"+id, "err", "invalid email format", "form", "update"), http.StatusSeeOther)
					return
				}
				var manufacturer _model.Manufacturer
				db.First(&manufacturer, id)
				if name == manufacturer.Name && phone == manufacturer.Phone && email == manufacturer.Email {
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer/"+id, "err", "no changes detected", "form", "update"), http.StatusSeeOther)
					return
				}
				manufacturer.Name = name
				manufacturer.Phone = phone
				manufacturer.Email = email
				db.Save(&manufacturer)
				http.Redirect(w, r, _util.URLBuilder("/app/manufacturer/"+id, "success", "manufacturer updated", "form", "update"), http.StatusSeeOther)
			},
			_middleware.Init, _middleware.ParseForm, _middleware.Auth,
		)
	})
}

func CreateEquipment(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /app/manufacturer/{id}", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				fmt.Println(r.Form)
				id := r.PathValue("id")
				name := r.Form.Get("nickname")
				number := r.Form.Get("number")
				submitRedirect := r.Form.Get("submitRedirect")
				redirectURL := _util.ConditionalString(
					submitRedirect == "", 
					"/app/manufacturer/"+id, 
					submitRedirect,
				)
				if name == "" || number == "" {
					http.Redirect(w, r, _util.URLBuilder(redirectURL, "form", "create", "err", "all fields required", "nickname", name, "serialNumber", number), http.StatusSeeOther)
					return
				}
				photo, _, err := r.FormFile("photo")
				if err != nil {
					http.Redirect(w, r, _util.URLBuilder(redirectURL, "form", "create", "err", "photo required", "nickname", name, "serialNumber", number), http.StatusSeeOther)
					return
				}
				defer photo.Close()
				if photo == nil {
					http.Redirect(w, r, _util.URLBuilder(redirectURL, "form", "create", "err", "photo required", "nickname", name, "serialNumber", number), http.StatusSeeOther)
					return
				}
				photoBytes, err := io.ReadAll(photo)
				if err != nil {
					http.Redirect(w, r, _util.URLBuilder(redirectURL, "form", "create", "err", "photo required", "nickname", name, "serialNumber", number), http.StatusSeeOther)
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
				http.Redirect(w, r, _util.URLBuilder(redirectURL, "success", "equipment created", "form", "create"), http.StatusSeeOther)
			},
			_middleware.Init, _middleware.ParseForm, _middleware.ParseMultipartForm, _middleware.Auth,
		)
	})
}

func UpdateEquipment(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /app/equipment/{id}/update", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
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
			_middleware.Init, _middleware.ParseMultipartForm, _middleware.Auth,
		)
	})
}

func DeleteEquipment(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /app/equipment/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
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
			_middleware.Init, _middleware.ParseForm, _middleware.Auth,
		)
	})
}

func CreateTicket(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /app/ticket/public", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				securityToken := r.Form.Get("publicSecurityToken")
				redirectURL := "/app/ticket/public"
				if securityToken == os.Getenv("ADMIN_REDIRECT_TOKEN") {
					redirectURL = "/app/ticket"
				}
				if redirectURL == "/app/ticket" {
					_middleware.Auth(customContext, w, r)
					securityToken = ""
				}
				creator := r.Form.Get("creator")
				item := r.Form.Get("item")
				problem := r.Form.Get("problem")
				location := r.Form.Get("location")
				if creator == "" || item == "" || problem == "" || location == "" {
					http.Redirect(w, r, _util.URLBuilder(redirectURL, "publicSecurityToken", securityToken, "err", "all fields required", "creator", creator, "item", item, "problem", problem, "location", location), http.StatusSeeOther)
					return
				}
				photo, _, err := r.FormFile("photo")
				if err != nil {
					http.Redirect(w, r, _util.URLBuilder(redirectURL, "publicSecurityToken", securityToken, "err", "photo required", "creator", creator, "item", item, "problem", problem, "location", location), http.StatusSeeOther)
					return
				}
				defer photo.Close()
				if photo == nil {
					http.Redirect(w, r, _util.URLBuilder(redirectURL, "publicSecurityToken", securityToken, "err", "photo required", "creator", creator, "item", item, "problem", problem, "location", location), http.StatusSeeOther)
					return
				}
				photoBytes, err := io.ReadAll(photo)
				if err != nil {
					http.Redirect(w, r, _util.URLBuilder(redirectURL, "publicSecurityToken", securityToken, "err", "photo required", "creator", creator, "item", item, "problem", problem, "location", location), http.StatusSeeOther)
					return
				}
				ticket := _model.Ticket{
					Creator:  creator,
					Item:     item,
					Problem:  problem,
					Location: location,
					Photo:    photoBytes,
					Priority: _model.TicketPriorityUnspecified,
					Status:   _model.TicketStatusNew,
					Notes:    "",
					Owner:    "",
				}
				db.Create(&ticket)
				http.Redirect(w, r, _util.URLBuilder(redirectURL, "publicSecurityToken", securityToken, "success", "your ticket has been created, thank you!"), http.StatusSeeOther)
			},
			_middleware.Init, _middleware.ParseMultipartForm,
		)
	})
}

func UpdateTicket(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /app/ticket/{id}/update", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				creator := r.Form.Get("creator")
				item := r.Form.Get("item")
				problem := r.Form.Get("problem")
				location := r.Form.Get("location")
				photo, _, err := r.FormFile("photo")
				defer func() {
					if photo != nil {
						photo.Close()
					}
				}()
				if err != nil && err != http.ErrMissingFile {
					http.Redirect(w, r, _util.URLBuilder("/app/ticket/"+id, "err", "failed to retrieve file"), http.StatusSeeOther)
					return
				}
				var ticket _model.Ticket
				db.First(&ticket, id)
				if photo != nil {
					photoBytes, err := io.ReadAll(photo)
					if err != nil {
						http.Error(w, "Failed to read file", http.StatusBadRequest)
						return
					}
					ticket.Photo = photoBytes
				}
				ticket.Creator = creator
				ticket.Item = item
				ticket.Problem = problem
				ticket.Location = location
				db.Save(&ticket)
				http.Redirect(w, r, "/app/ticket/"+id, http.StatusSeeOther)
			},
			_middleware.Init, _middleware.ParseMultipartForm, _middleware.Auth,
		)
	})
}

func DeleteTicket(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /app/ticket/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				keyword := strings.ToLower(r.Form.Get("keyword"))
				if keyword != "yes" {
					http.Redirect(w, r, _util.URLBuilder("/app/ticket/"+id, "err", "invalid delete keyword provided"), http.StatusSeeOther)
					return
				}
				var ticket _model.Ticket
				db.First(&ticket, id)
				db.Delete(&ticket, id)
				http.Redirect(w, r, "/app/ticket", http.StatusSeeOther)
			},
			_middleware.Init, _middleware.ParseForm, _middleware.Auth,
		)
	})
}

func UpdateTicketPublicDetails(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /app/ticket/{id}/publicdetails", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				owner := r.Form.Get("owner")
				priority := r.Form.Get("priority")
				status := r.Form.Get("status")
				notes := r.Form.Get("notes")
				var ticket _model.Ticket
				db.First(&ticket, id)
				ticket.Owner = owner
				ticket.Priority = _model.TicketPriority(priority) // Convert priority to _model.TicketPriority type
				ticket.Status = _model.TicketStatus(status)       // Convert status to _model.TicketStatus type
				ticket.Notes = notes
				db.Save(&ticket)
				http.Redirect(w, r, _util.URLBuilder("/app/ticket/"+id, "publicSecurityToken", r.Form.Get("publicSecurityToken")), http.StatusSeeOther)
			},
			_middleware.Init, _middleware.ParseMultipartForm,
		)
	})
}

func AssignTicket(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /app/ticket/{id}/assign", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				equipmentID := r.Form.Get("equipmentID")
				if equipmentID == "" {
					http.Redirect(w, r, _util.URLBuilder("/app/ticket/"+id, "err", "equipment selection required"), http.StatusSeeOther)
					return
				}
				// convert equipmentID to uint
				uintEquipmentID, err := _util.ConvertStringToUint(equipmentID)
				if err != nil {
					http.Redirect(w, r, _util.URLBuilder("/app/ticket/"+id, "err", "invalid equipment selection"), http.StatusSeeOther)
					return
				}
				var ticket _model.Ticket
				db.First(&ticket, id)
				ticket.EquipmentID = &uintEquipmentID
				ticket.Status = _model.TicketStatusActive
				db.Save(&ticket)
				http.Redirect(w, r, "/app/ticket/"+id, http.StatusSeeOther)
			},
			_middleware.Init, _middleware.ParseForm, _middleware.Auth,
		)
	})
}

func TicketResetEquipment(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/ticket/{id}/resetEquipment", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				var ticket _model.Ticket
				db.First(&ticket, id)
				ticket.EquipmentID = nil
				ticket.Status = _model.TicketStatusNew
				db.Save(&ticket)
				http.Redirect(w, r, "/app/ticket/"+id, http.StatusSeeOther)
			},
			_middleware.Init, _middleware.Auth,
		)
	})
}
