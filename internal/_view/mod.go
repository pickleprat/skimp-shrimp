package _view

import (
	"cfasuite/internal/_components"
	"cfasuite/internal/_middleware"
	"cfasuite/internal/_model"
	"cfasuite/internal/_util"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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
				[x-cloak] { 
					display: none !important; 
				}
				
				.flex-indicator {
					display:none;
				}
								
				.htmx-request .flex-indicator {
					display: flex;
				}
				.htmx-request.flex-indicator {
					display: flex;
				}
			</style>
			<script>
				function qsa(sel) {
					return document.querySelectorAll(sel);
				}
				function qs(sel) {
					return document.querySelector(sel);
				}
			</script>
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

func Home(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				b := NewViewBuilder("Repairs Log - Login", []string{
					_components.Banner(false, ""),
					_components.MainLoader(),
					_components.Root(
						_components.CenterContentWrapper(
							_components.LoginForm(r.URL.Query().Get("err"), r.URL.Query().Get("username"), r.URL.Query().Get("password")),
						),
					),
					_components.Footer(),
				})
				w.Write(b.Build())
			},
			_middleware.Init,
		)
	})
}

func CreateManufacturers(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/manufacturer/create", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				var manufacturers []_model.Manufacturer
				submitRedirect := r.URL.Query().Get("submitRedirect")
				db.Find(&manufacturers)
				b := NewViewBuilder("Repairs Log - Create Manufacturers", []string{
					_components.Banner(true, _components.AppNavMenu(r.URL.Path)),
					_components.MainLoader(),
					_components.Root(
						_components.CenterContentWrapper(
							_components.SimpleTopNav(
								_components.LinkButton("/app/ticket/view", "View Tickets", false),
								_components.LinkButton("/app/ticket/create", "Create Tickets", false),
								_components.LinkButton("/app/manufacturer/view", "All Manufacturers", false),
								_components.LinkButton("/app/manufacturer/create", "Create Manufacturers", true),
							),
							_components.CreateManufacturerForm(r.URL.Query().Get("err"), r.URL.Query().Get("success"), r.URL.Query().Get("name"), r.URL.Query().Get("phone"), r.URL.Query().Get("email"), submitRedirect),
						),
					),
					_components.Footer(),
				})
				w.Write(b.Build())
			},
			_middleware.Init, _middleware.Auth,
		)
	})
}

func ViewManufacturers(mux *http.ServeMux, db *gorm.DB) {
    mux.HandleFunc("GET /app/manufacturer/view", func(w http.ResponseWriter, r *http.Request) {
        _middleware.MiddlewareChain(w, r,
            func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
                var manufacturers []_model.Manufacturer
                // Sorting manufacturers alphabetically by name (case-insensitive)
                db.Order("LOWER(name)").Find(&manufacturers)
                b := NewViewBuilder("Repairs Log - View Manufacturers", []string{
                    _components.Banner(true, _components.AppNavMenu(r.URL.Path)),
                    _components.MainLoader(),
                    _components.Root(
                        _components.CenterContentWrapper(
                            _components.SimpleTopNav(
                                _components.LinkButton("/app/ticket/view", "View Tickets", false),
                                _components.LinkButton("/app/ticket/create", "Create Tickets", false),
                                _components.LinkButton("/app/manufacturer/view", "All Manufacturers", true),
                                _components.LinkButton("/app/manufacturer/create", "Create Manufacturers", false),
                            ),
                            _components.ManufacturerList(manufacturers, r.URL.Query().Get("err")),
                        ),
                    ),
                    _components.Footer(),
                })
                w.Write(b.Build())
            },
            _middleware.Init, _middleware.Auth,
        )
    })
}

func ViewManufacturer(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/manufacturer/{id}/update", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				var manufacturer _model.Manufacturer
				db.First(&manufacturer, id)
				b := NewViewBuilder("Repairs Log - " + manufacturer.Name, []string{
					_components.Banner(true, _components.AppNavMenu(r.URL.Path)),
					_components.MainLoader(),
					_components.Root(
						_components.CenterContentWrapper(
							_components.SimpleTopNav(
								_components.LinkButton("/app/manufacturer/view", "Back", false),
								_components.LinkButton("/app/manufacturer/"+ id +"/update", "Update", true),
                                _components.LinkButton("/app/manufacturer/"+ id +"/delete", "Delete", false),
								_components.LinkButton(fmt.Sprintf("/app/manufacturer/%d/equipment/create", manufacturer.ID), "Create Equipment", false),
								_components.LinkButton(fmt.Sprintf("/app/manufacturer/%d/equipment/view", manufacturer.ID), "View Equipment", false),
                            ),
							_components.UpdateManufacturerForm(manufacturer, r.URL.Query().Get("err"), r.URL.Query().Get("success")),
						),
					),
					_components.Footer(),
				})
				w.Write(b.Build())
			},
			_middleware.Init, _middleware.Auth,
		)
	})
}

func DeleteManufacturer(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/manufacturer/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				var manufacturer _model.Manufacturer
				db.First(&manufacturer, id)
				b := NewViewBuilder("Repairs Log - " + manufacturer.Name, []string{
					_components.Banner(true, _components.AppNavMenu(r.URL.Path)),
					_components.MainLoader(),
					_components.Root(
						_components.CenterContentWrapper(
							_components.SimpleTopNav(
								_components.LinkButton("/app/manufacturer/view", "Back", false),
								_components.LinkButton("/app/manufacturer/"+ id +"/update", "Update", false),
								_components.LinkButton("/app/manufacturer/"+ id +"/delete", "Delete", true),
								_components.LinkButton(fmt.Sprintf("/app/manufacturer/%d/equipment/create", manufacturer.ID), "Create Equipment", false),
								_components.LinkButton(fmt.Sprintf("/app/manufacturer/%d/equipment/view", manufacturer.ID), "View Equipment", false),
							),
							_components.DeleteManufacturerForm(manufacturer, r.URL.Query().Get("err")),
						),
					),
					_components.Footer(),
				})
				w.Write(b.Build())
			},
			_middleware.Init, _middleware.Auth,
		)
	})
}

func CreateEquipment(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/manufacturer/{id}/equipment/create", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				var manufacturer _model.Manufacturer
				db.First(&manufacturer, id)
				b := NewViewBuilder("Repairs Log - Create Equipment", []string{
					_components.Banner(true, _components.AppNavMenu(r.URL.Path)),
					_components.MainLoader(),
					_components.Root(
						_components.CenterContentWrapper(
							_components.SimpleTopNav(
                                _components.LinkButton("/app/manufacturer/view", "Back", false),
                                _components.LinkButton("/app/manufacturer/"+ id +"/update", "Update", false),
                                _components.LinkButton("/app/manufacturer/"+ id +"/delete", "Delete", false),
								_components.LinkButton(fmt.Sprintf("/app/manufacturer/%d/equipment/create", manufacturer.ID), "Create Equipment", true),
								_components.LinkButton(fmt.Sprintf("/app/manufacturer/%d/equipment/view", manufacturer.ID), "View Equipment", false),
                            ),
							_components.CreateEquipmentForm(
								manufacturer, 
								r.URL.Query().Get("err"), 
								r.URL.Query().Get("success"),
								r.URL.Query().Get("submitRedirect"),
								r.URL.Query().Get("nickname"),
								r.URL.Query().Get("serialNumber"),
							),
						),
					),
					_components.Footer(),
				})
				w.Write(b.Build())
			},
			_middleware.Init, _middleware.Auth,
		)
	})
}

func ViewEquipment(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/manufacturer/{id}/equipment/view", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				var manufacturer _model.Manufacturer
				var equipment []_model.Equipment
				db.First(&manufacturer, id)
				db.Where("manufacturer_id = ?", manufacturer.ID).Find(&equipment)
				b := NewViewBuilder("Repairs Log - View Equipment", []string{
					_components.Banner(true, _components.AppNavMenu(r.URL.Path)),
					_components.MainLoader(),
					_components.Root(
						_components.CenterContentWrapper(
							_components.SimpleTopNav(
								_components.LinkButton("/app/manufacturer/view", "Back", false),
								_components.LinkButton("/app/manufacturer/"+ id +"/update", "Update", false),
                                _components.LinkButton("/app/manufacturer/"+ id +"/delete", "Delete", false),
								_components.LinkButton(fmt.Sprintf("/app/manufacturer/%d/equipment/create", manufacturer.ID), "Create Equipment", false),
								_components.LinkButton(fmt.Sprintf("/app/manufacturer/%d/equipment/view", manufacturer.ID), "View Equipment", true),
							),
							_components.EquipmentList(equipment),
						),
					),
					_components.Footer(),
				})
				w.Write(b.Build())
			},
			_middleware.Init, _middleware.Auth,
		)
	})
}

func Equipment(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/equipment/{id}", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				form := r.URL.Query().Get("form")
				if form == "" {
					form = "update"
				}
				var equipment _model.Equipment
				var manufacturer _model.Manufacturer
				db.First(&equipment, id)
				db.First(&manufacturer, equipment.ManufacturerID)
				b := NewViewBuilder("Repairs Log - App", []string{
					_components.Banner(true, _components.AppNavMenu(r.URL.Path)),
					_components.MainLoader(),
					_components.Root(
						_components.CenterContentWrapper(
							_components.EquipmentDetails(equipment, manufacturer, form),
							_util.ConditionalString(
								form == "update",
								_components.UpdateEquipmentForm(equipment, r.URL.Query().Get("err"), r.URL.Query().Get("success")),
								"",
							),
							_util.ConditionalString(
								form == "delete",
								_components.DeleteEquipmentForm(equipment, r.URL.Query().Get("err"), r.URL.Query().Get("success")),
								"",
							),
						),
					),
					_components.Footer(),
				})
				w.Write(b.Build())
			},
			_middleware.Init, _middleware.Auth,
		)
	})
}

func GetEquipmentQRCode(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/equipment/{id}/downloadticketqr", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
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
			_middleware.Init, _middleware.Auth, _middleware.IncludePNG,
		)
	})
}

func EquipmentTicket(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/equipment/{id}/ticketform", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
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
					_components.MainLoader(),
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
			_middleware.Init,
		)
	})
}

func TicketForm(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/ticket/public", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				token := r.URL.Query().Get("publicSecurityToken")
				if token != os.Getenv("PUBLIC_SECURITY_TOKEN") {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
				b := NewViewBuilder("Repairs Log - Public Tickets", []string{
					_components.Banner(false, ""),
					_components.MainLoader(),
					_components.Root(
						_components.CenterContentWrapper(
							_components.CreateTicketForm(r, token, "/app/ticket/public"),
						),
					),
					_components.Footer(),
				})
				w.Write(b.Build())
			},
			_middleware.Init,
		)
	})
}

func AdminViewTicket(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/ticket/{id}/view", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				var ticket _model.Ticket
				id := r.PathValue("id")
				form := r.URL.Query().Get("form")
				if err := db.First(&ticket, id).Error; err != nil {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				var manufacturers []_model.Manufacturer
				db.Find(&manufacturers)
				if ticket.EquipmentID == nil && form == "" {
					form = "assign"
				}
				if ticket.EquipmentID != nil && form == "" {
					form = "public"
				}
				b := NewViewBuilder("Repairs Log - Tickets", []string{
					_components.Banner(true, _components.AppNavMenu(r.URL.Path)),
					_components.MainLoader(),
					_components.Root(
						_components.CenterContentWrapper(
							_components.TicketDetails(ticket, r.URL.Query().Get("err")),
							_components.TicketSettings(ticket, form),
							_util.ConditionalString(
								form == "update",
								_components.UpdateTicketForm(ticket, r.URL.Query().Get("err"), r.URL.Query().Get("success")),
								"",
							),
							_util.ConditionalString(
								form == "delete",
								_components.DeleteTicketForm(ticket, r.URL.Query().Get("err")),
								"",
							),
							_util.ConditionalString(
								form == "assign",
								_components.TicketAssignmentForm(ticket, manufacturers, db),
								"",
							),
							_util.ConditionalString(
								form == "public",
								_components.TicketPublicDetailsForm(ticket, r.URL.Query().Get("err"), r.URL.Query().Get("success")),
								"",
							),
						),
					),
					_components.Footer(),
				})
				w.Write(b.Build())
			},
			_middleware.Init, _middleware.Auth,
		)
	})
}

func AdminCreateTickets(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/ticket/create", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				var manufacturers []_model.Manufacturer
				db.Find(&manufacturers)
				b := NewViewBuilder("Repairs Log - Create Tickets", []string{
					_components.Banner(true, _components.AppNavMenu(r.URL.Path)),
					_components.MainLoader(),
					_components.Root(
						_components.CenterContentWrapper(
							_components.SimpleTopNav(
								_components.LinkButton("/app/ticket/view", "View Tickets", false),
								_components.LinkButton("/app/ticket/create", "Create Tickets", true),
								_components.LinkButton("/app/manufacturer/view", "All Manufacturers", false),
								_components.LinkButton("/app/manufacturer/create", "Create Manufacturers", false),
							),
							_components.CreateTicketForm(r, "", "/form/ticket/admin"),
						),
					),
					_components.Footer(),
				})
				w.Write(b.Build())
			},
			_middleware.Init, _middleware.Auth,
		)
	})
}

func AdminViewTickets(mux *http.ServeMux, db *gorm.DB) {
    mux.HandleFunc("GET /app/ticket/view", func(w http.ResponseWriter, r *http.Request) {
        _middleware.MiddlewareChain(w, r,
            func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
                var tickets []_model.Ticket
                if err := db.Order("created_at DESC").Find(&tickets).Error; err != nil {
                    http.Error(w, err.Error(), http.StatusInternalServerError)
                    return
                }
                b := NewViewBuilder("Repairs Log - View Tickets", []string{
                    _components.Banner(true, _components.AppNavMenu(r.URL.Path)),
                    _components.MainLoader(),
					_components.Root(
						_components.CenterContentWrapper(
							_components.SimpleTopNav(
								_components.LinkButton("/app/ticket/view", "View Tickets", true),
								_components.LinkButton("/app/ticket/create", "Create Tickets", false),
								_components.LinkButton("/app/manufacturer/view", "All Manufacturers", false),
								_components.LinkButton("/app/manufacturer/create", "Create Manufacturers", false),
							),
							_components.TicketList(tickets),
						),
					),
                    _components.Footer(),
                })
                w.Write(b.Build())
            },
            _middleware.Init, _middleware.Auth,
        )
    })
}

