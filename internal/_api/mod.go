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
					http.Redirect(w, r, "/app/ticket/view", http.StatusSeeOther)
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
	mux.HandleFunc("POST /form/manufacturer/create", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				name := r.Form.Get("name")
				phone := r.Form.Get("phone")
				email := r.Form.Get("email")
				submitRedirect := r.Form.Get("submitRedirect")
				if name == "" || phone == "" || email == "" {
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer/create", "err", "all fields required", "name", name, "phone", phone, "email", email, "submitRedirect", submitRedirect), http.StatusSeeOther)
					return
				}
				if _util.IsValidPhoneNumber(phone) == false {
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer/create", "err", "invalid phone format", "name", name, "phone", phone, "email", email, "submitRedirect", submitRedirect), http.StatusSeeOther)
					return
				}
				if _util.IsValidEmail(email) == false {
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer/create", "err", "invalid email format", "name", name, "phone", phone, "email", email, "submitRedirect", submitRedirect), http.StatusSeeOther)
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
					_util.URLBuilder("/app/manufacturer/create", "success", "manufacturer created"), 
					_util.URLBuilder(submitRedirect, "success", "manufacturer created"),
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
	mux.HandleFunc("POST /form/manufacturer/{id}/equipment/create", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				name := r.Form.Get("nickname")
				number := r.Form.Get("number")
				submitRedirect := r.Form.Get("submitRedirect")
				redirectURL := _util.ConditionalString(
					submitRedirect == "", 
					fmt.Sprintf("/app/manufacturer/%s/equipment/create", id), 
					submitRedirect,
				)
				if name == "" || number == "" {
					http.Redirect(w, r, _util.URLBuilder(redirectURL, "err", "all fields required", "nickname", name, "serialNumber", number, "submitRedirect", submitRedirect), http.StatusSeeOther)
					return
				}
				photo, _, err := r.FormFile("photo")
				if err != nil {
					http.Redirect(w, r, _util.URLBuilder(redirectURL, "err", "photo required", "nickname", name, "serialNumber", number, "submitRedirect", submitRedirect), http.StatusSeeOther)
					return
				}
				defer photo.Close()
				if photo == nil {
					http.Redirect(w, r, _util.URLBuilder(redirectURL, "err", "photo required", "nickname", name, "serialNumber", number, "submitRedirect", submitRedirect), http.StatusSeeOther)
					return
				}
				photoBytes, err := io.ReadAll(photo)
				if err != nil {
					http.Redirect(w, r, _util.URLBuilder(redirectURL, "err", "photo required", "nickname", name, "serialNumber", number, "submitRedirect", submitRedirect), http.StatusSeeOther)
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
				http.Redirect(w, r, _util.URLBuilder(redirectURL, "success", "equipment created"), http.StatusSeeOther)
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
				photo, _, err := r.FormFile("photo")
				if name == "" || number == "" {
					http.Redirect(w, r, _util.URLBuilder("/app/equipment/"+id, "err", "all fields required"), http.StatusSeeOther)
					return
				}
				defer func() {
					if photo != nil {
						photo.Close()
					}
				}()
				if err != nil && err != http.ErrMissingFile {
					http.Redirect(w, r, _util.URLBuilder("/app/equipment/"+id, "err", "failed to retrieve file"), http.StatusSeeOther)
					return
				}

				var equipment _model.Equipment
				db.First(&equipment, id)
				if name == equipment.Nickname && number == equipment.SerialNumber && photo == nil {
					http.Redirect(w, r, _util.URLBuilder("/app/equipment/"+id, "err", "no changes detected"), http.StatusSeeOther)
					return
				}
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
				http.Redirect(w, r, _util.URLBuilder("/app/equipment/"+id, "success", "equipment updated"), http.StatusSeeOther)
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
					http.Redirect(w, r, _util.URLBuilder(fmt.Sprintf("/app/equipment/%d", equipment.ID), "err", "invalid delete name provided", "form", "delete"), http.StatusSeeOther)
					return
				}
				db.Delete(&equipment, id)
				http.Redirect(w, r, fmt.Sprintf(_util.URLBuilder("/app/manufacturer/%d", "success", "equipment deleted"), manufacturer.ID), http.StatusSeeOther)
			},
			_middleware.Init, _middleware.ParseForm, _middleware.Auth,
		)
	})
}

func CreateTicketPublic(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /app/ticket/public", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				securityToken := r.Form.Get("publicSecurityToken")
				if securityToken != os.Getenv("PUBLIC_SECURITY_TOKEN") {
					http.Error(w, "invalid security token", http.StatusUnauthorized)
					return
				}
				creator := r.Form.Get("creator")
				item := r.Form.Get("item")
				problem := r.Form.Get("problem")
				location := r.Form.Get("location")
				if creator == "" || item == "" || problem == "" || location == "" {
					http.Redirect(w, r, _util.URLBuilder("/app/ticket/public", "publicSecurityToken", securityToken, "err", "all fields required", "creator", creator, "item", item, "problem", problem, "location", location), http.StatusSeeOther)
					return
				}
				ticket := _model.Ticket{
					Creator:  creator,
					Item:     item,
					Problem:  problem,
					Location: location,
					Priority: _model.TicketPriorityUnspecified,
					Status:   _model.TicketStatusNew,
					Notes:    "",
					Owner:    "",
				}
				db.Create(&ticket)
				http.Redirect(w, r, _util.URLBuilder("/app/ticket/public", "publicSecurityToken", securityToken, "success", "your ticket has been created, thank you!"), http.StatusSeeOther)
			},
			_middleware.Init, _middleware.ParseForm,
		)
	})

}

func CreateTicketAdmin(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /form/ticket/admin", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				creator := r.Form.Get("creator")
				item := r.Form.Get("item")
				problem := r.Form.Get("problem")
				location := r.Form.Get("location")
				if creator == "" || item == "" || problem == "" || location == "" {
					http.Redirect(w, r, _util.URLBuilder("/app/ticket/create", "err", "all fields required", "creator", creator, "item", item, "problem", problem, "location", location), http.StatusSeeOther)
					return
				}
				ticket := _model.Ticket{
					Creator:  creator,
					Item:     item,
					Problem:  problem,
					Location: location,
					Priority: _model.TicketPriorityUnspecified,
					Status:   _model.TicketStatusNew,
					Notes:    "",
					Owner:    "",
				}
				db.Create(&ticket)
				http.Redirect(w, r, _util.URLBuilder("/app/ticket/create", "success", "your ticket has been created, thank you!"), http.StatusSeeOther)
			},
			_middleware.Init, _middleware.ParseMultipartForm, _middleware.Auth,
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
				if creator == "" || item == "" || problem == "" || location == "" {
					http.Redirect(w, r, _util.URLBuilder("/app/ticket/"+id, "err", "all fields required", "form", "update"), http.StatusSeeOther)
					return
				}
				var ticket _model.Ticket
				db.First(&ticket, id)
				if creator == ticket.Creator && item == ticket.Item && problem == ticket.Problem && location == ticket.Location {
					http.Redirect(w, r, _util.URLBuilder("/app/ticket/"+id, "err", "no changes detected", "form", "update"), http.StatusSeeOther)
					return
				}
				ticket.Creator = creator
				ticket.Item = item
				ticket.Problem = problem
				ticket.Location = location
				db.Save(&ticket)
				http.Redirect(w, r, _util.URLBuilder("/app/ticket/"+id, "form", "update", "success", "ticket updated"), http.StatusSeeOther)
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
					http.Redirect(w, r, _util.URLBuilder("/app/ticket/"+id, "form", "delete", "err", "invalid delete keyword provided"), http.StatusSeeOther)
					return
				}
				var ticket _model.Ticket
				db.First(&ticket, id)
				db.Delete(&ticket, id)
				http.Redirect(w, r, _util.URLBuilder("/app/ticket", "success", "ticket deleted"), http.StatusSeeOther)
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
				if owner == ticket.Owner && priority == string(ticket.Priority) && status == string(ticket.Status) && notes == ticket.Notes {
					http.Redirect(w, r, _util.URLBuilder("/app/ticket/"+id, "form", "public", "err", "no changes detected"), http.StatusSeeOther)
					return
				}
				ticket.Owner = owner
				ticket.Priority = _model.TicketPriority(priority) // Convert priority to _model.TicketPriority type
				ticket.Status = _model.TicketStatus(status)       // Convert status to _model.TicketStatus type
				ticket.Notes = notes
				db.Save(&ticket)
				http.Redirect(w, r, _util.URLBuilder("/app/ticket/"+id, "form", "public", "success", "public details updated"), http.StatusSeeOther)
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
