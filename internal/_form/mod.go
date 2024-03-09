package _form

import (
	"cfasuite/internal/_middleware"
	"cfasuite/internal/_model"
	"cfasuite/internal/_util"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

func Login(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /form/login", func(w http.ResponseWriter, r *http.Request) {
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
						Path: "/",
					})
					http.Redirect(w, r, "/app", http.StatusSeeOther)
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
				if name == ""  {
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer", "err", "name is required", "name", name, "phone", phone, "email", email), http.StatusSeeOther)
					return
				}
				manufacturer := _model.Manufacturer{
					Name:  name,
					Phone: phone,
					Email: email,
				}
				db.Create(&manufacturer)
				http.Redirect(w, r, _util.URLBuilder("/app/manufacturer", "success", "manufacturer created"), http.StatusSeeOther)
			},
			_middleware.Init, _middleware.Auth, _middleware.ParseForm,
		)
	})
}

func DeleteManufacturer(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /form/manufacturer/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				name := strings.ToLower(r.Form.Get("name"))
				var manufacturer _model.Manufacturer
				db.First(&manufacturer, id)
				if name != strings.ToLower(manufacturer.Name) {
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer/"+id+"/delete", "err", "invalid delete name provided", "form", "delete"), http.StatusSeeOther)
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
	mux.HandleFunc("POST /form/manufacturer/{id}/update", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				name := r.Form.Get("name")
				phone := r.Form.Get("phone")
				email := r.Form.Get("email")
				if name == "" || phone == "" || email == "" {
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer/"+id+"/update", "err", "all fields required", "form", "update"), http.StatusSeeOther)
					return
				}
				if _util.IsValidPhoneNumber(phone) == false {
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer/"+id+"/update", "err", "invalid phone format", "form", "update"), http.StatusSeeOther)
					return
				}
				if _util.IsValidEmail(email) == false {
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer/"+id+"/update", "err", "invalid email format", "form", "update"), http.StatusSeeOther)
					return
				}
				var manufacturer _model.Manufacturer
				db.First(&manufacturer, id)
				if name == manufacturer.Name && phone == manufacturer.Phone && email == manufacturer.Email {
					http.Redirect(w, r, _util.URLBuilder("/app/manufacturer/"+id+"/update", "err", "no changes detected", "form", "update"), http.StatusSeeOther)
					return
				}
				manufacturer.Name = name
				manufacturer.Phone = phone
				manufacturer.Email = email
				db.Save(&manufacturer)
				http.Redirect(w, r, _util.URLBuilder("/app/manufacturer/"+id+"/update", "success", "manufacturer updated", "form", "update"), http.StatusSeeOther)
			},
			_middleware.Init, _middleware.ParseForm, _middleware.Auth,
		)
	})
}

func CreateEquipment(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /form/manufacturer/{id}/createEquipment", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				name := r.Form.Get("nickname")
				serialNumber := r.Form.Get("serialNumber")
				modelNumber := r.Form.Get("modelNumber")
				redirectURL := "/app/manufacturer/" + id 
				if name == "" || serialNumber == "" || modelNumber == "" {
					http.Redirect(w, r, _util.URLBuilder(redirectURL, "err", "all fields required", "nickname", name, "serialNumber", serialNumber, "modelNumber", modelNumber), http.StatusSeeOther)
					return
				}
				photo, _, err := r.FormFile("photo")
				if err != nil {
					http.Redirect(w, r, _util.URLBuilder(redirectURL, "err", "photo required", "nickname", name, "serialNumber", serialNumber, "modelNumber", modelNumber), http.StatusSeeOther)
					return
				}
				defer photo.Close()
				if photo == nil {
					http.Redirect(w, r, _util.URLBuilder(redirectURL, "err", "photo required", "nickname", name, "serialNumber", serialNumber, "modelNumber", modelNumber), http.StatusSeeOther)
					return
				}
				photoBytes, err := io.ReadAll(photo)
				if err != nil {
					http.Redirect(w, r, _util.URLBuilder(redirectURL, "err", "photo required", "nickname", name, "serialNumber", serialNumber, "modelNumber", modelNumber), http.StatusSeeOther)
					return
				}
				var manufacturer _model.Manufacturer
				db.First(&manufacturer, id)
				equipment := _model.Equipment{
					Nickname:       name,
					SerialNumber:   serialNumber,
					ModelNumber:    modelNumber,
					Archived: 	 	false,
					Photo:          photoBytes,
					ManufacturerID: manufacturer.ID,
				}
				db.Create(&equipment)
				http.Redirect(w, r, _util.URLBuilder(redirectURL, "success", "equipment created"), http.StatusSeeOther)
			},
			_middleware.Init, _middleware.ParseForm, _middleware.ParseMultipartForm, _middleware.Auth,
		)
	})
}

func UpdateEquipment(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /form/equipment/{id}/update", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				name := r.Form.Get("nickname")
				serialNumber := r.Form.Get("serialNumber")
				modelNumber := r.Form.Get("modelNumber")
				archived := r.Form.Get("archived")
				photo, _, err := r.FormFile("photo")
				if name == "" || serialNumber == "" || modelNumber == "" {
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
				var archivedString string
				if equipment.Archived {
					archivedString = "true"
				} else {
					archivedString = "false"
				}
				if name == equipment.Nickname && serialNumber == equipment.SerialNumber && modelNumber == equipment.ModelNumber && archived == archivedString && photo == nil {
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
				equipment.SerialNumber = serialNumber
				equipment.ModelNumber = modelNumber
				equipment.Archived = archived == "true"
				db.Save(&equipment)
				http.Redirect(w, r, _util.URLBuilder("/app/equipment/"+id, "success", "equipment updated"), http.StatusSeeOther)
			},
			_middleware.Init, _middleware.ParseMultipartForm, _middleware.Auth,
		)
	})
}

func DeleteEquipment(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /form/equipment/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				name := strings.ToLower(r.Form.Get("name"))
				var equipment _model.Equipment
				db.First(&equipment, id)
				var manufacturer _model.Manufacturer
				db.First(&manufacturer, equipment.ManufacturerID)
				if name != strings.ToLower(equipment.Nickname) {
					http.Redirect(w, r, _util.URLBuilder(fmt.Sprintf("/app/equipment/%d/delete", equipment.ID), "err", "invalid delete name provided", "form", "delete"), http.StatusSeeOther)
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
	mux.HandleFunc("POST /form/ticket/public", func(w http.ResponseWriter, r *http.Request) {
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
				englishItem, err := _util.TranslateToEnglish(item)
				if err != nil {
					englishItem = item
				}
				englishProblem, err := _util.TranslateToEnglish(problem)
				if err != nil {
					englishProblem = problem
				}
				ticket := _model.Ticket{
					Creator:  creator,
					Item:     englishItem,
					Problem:  englishProblem,
					Location: location,
					Priority: _model.TicketPriorityLow,
					Status:   _model.TicketStatusActive,
					Notes:    "",
					Owner:    "",
				}
				db.Create(&ticket)
				http.Redirect(w, r, _util.URLBuilder("/app/ticket/public", "publicSecurityToken", securityToken, "success", "your ticket is created and is up for review | tu boleto ha sido creado y está listo para su revisión "), http.StatusSeeOther)
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
					http.Redirect(w, r, _util.URLBuilder("/app/ticket", "err", "all fields required", "creator", creator, "item", item, "problem", problem, "location", location), http.StatusSeeOther)
					return
				}
				ticket := _model.Ticket{
					Creator:  creator,
					Item:     item,
					Problem:  problem,
					Location: location,
					Priority: _model.TicketPriorityLow,
					Status:   _model.TicketStatusActive,
					Notes:    "",
					Owner:    "",
				}
				db.Create(&ticket)
				http.Redirect(w, r, _util.URLBuilder("/app/ticket", "success", "your ticket has been created, thank you!"), http.StatusSeeOther)
			},
			_middleware.Init, _middleware.ParseMultipartForm, _middleware.Auth,
		)
	})
}

func UpdateTicket(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /form/ticket/{id}/update", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				creator := r.Form.Get("creator")
				item := r.Form.Get("item")
				problem := r.Form.Get("problem")
				location := r.Form.Get("location")
				priority := r.Form.Get("priority")
				status := r.Form.Get("status")
				notes := r.Form.Get("notes")
				owner := r.Form.Get("owner")
				if creator == "" || item == "" || problem == "" || location == "" || priority == "" || status == "" {
					http.Redirect(w, r, _util.URLBuilder("/app/ticket/"+id, "err", "all fields required"), http.StatusSeeOther)
					return
				}
				var ticket _model.Ticket
				db.First(&ticket, id)
				ticket.Creator = creator
				ticket.Item = item
				ticket.Problem = problem
				ticket.Location = location
				ticket.Priority = _model.TicketPriority(priority)
				ticket.Status = _model.TicketStatus(status)
				ticket.Notes = notes
				ticket.Owner = owner
				db.Save(&ticket)
				http.Redirect(w, r, _util.URLBuilder("/app/ticket/"+id+"/update", "success", "ticket updated"), http.StatusSeeOther)
			},
			_middleware.Init, _middleware.ParseForm, _middleware.Auth,
		)
	})
}

func DeleteTicket(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /form/ticket/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
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


func AssignTicket(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /form/ticket/{id}/assign", func(w http.ResponseWriter, r *http.Request) {
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
	mux.HandleFunc("GET /form/ticket/{id}/resetEquipment", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				var ticket _model.Ticket
				db.First(&ticket, id)
				ticket.EquipmentID = nil
				ticket.Status = _model.TicketStatusActive
				db.Save(&ticket)
				http.Redirect(w, r, "/app/ticket/"+id, http.StatusSeeOther)
			},
			_middleware.Init, _middleware.Auth,
		)
	})
}

func CompleteTicket(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("POST /form/ticket/{id}/complete", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				cost := r.Form.Get("cost")
				repairNotes := r.Form.Get("repairNotes")
				if cost == "" {
					cost = "0"
				}
				costFloat, err := strconv.ParseFloat(cost, 64)
				if err != nil {
					http.Redirect(w, r, _util.URLBuilder("/app/ticket/"+id+"/complete", "err", "cost must be a valid number with no symbols or empty"), http.StatusSeeOther)
					return
				}
				if repairNotes == "" {
					http.Redirect(w, r, _util.URLBuilder("/app/ticket/"+id+"/complete", "err", "repair notes required"), http.StatusSeeOther)
					return
				}
				var ticket _model.Ticket
				db.First(&ticket, id)
				ticket.Status = _model.TicketStatusComplete
				ticket.Cost = costFloat
				ticket.RepairNotes = repairNotes
				db.Save(&ticket)
				http.Redirect(w, r, "/app/ticket", http.StatusSeeOther)
			},
			_middleware.Init, _middleware.ParseForm, _middleware.Auth,
		)
	})
}