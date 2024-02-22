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
				[x-cloak] { display: none !important; }
				.htmx-indicator{
					opacity:0;
					transition: opacity 500ms ease-in;
				}
				.htmx-request .htmx-indicator{
					opacity:1
				}
				.htmx-request.htmx-indicator{
					opacity:1
				}
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

func ManufacturerForm(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/manufacturer", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
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
			_middleware.Init, _middleware.Auth,
		)
	})
}



func Manufacturer(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/manufacturer/{id}", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
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
			_middleware.Init, _middleware.Auth,
		)
	})
}

func Equipment(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/equipment/{id}", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
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
					_components.Root(
						_components.CenterContentWrapper(
							_components.CreateTicketForm(r, token),
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

func Tickets(mux *http.ServeMux, db *gorm.DB) {
    mux.HandleFunc("GET /app/ticket", func(w http.ResponseWriter, r *http.Request) {
        _middleware.MiddlewareChain(w, r,
            func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
                var tickets []_model.Ticket
                ticketFilter := r.URL.Query().Get("ticketFilter")
                
                switch ticketFilter {
                case "", "needsAttention":
                    // Filter tickets where any of the pointer fields are empty
                    if err := db.Where("priority IS NULL OR owner IS NULL OR status IS NULL OR notes IS NULL").Find(&tickets).Error; err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                    }
                case "all":
                    // Retrieve all tickets
                    if err := db.Find(&tickets).Error; err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                    }
                case "completed":
                    // Retrieve completed tickets
                    if err := db.Where("status = ?", "completed").Find(&tickets).Error; err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                    }
                default:
                    // Handle invalid ticketFilter value
                    http.Error(w, "Invalid ticket filter", http.StatusBadRequest)
                    return
                }

                b := NewViewBuilder("Repairs Log - Tickets", []string{
                    _components.Banner(true, _components.AppNavMenu(r.URL.Path)),
                    _components.Root(
                        _components.CenterContentWrapper(
                            _components.CreateTicketForm(r, os.Getenv("ADMIN_REDIRECT_TOKEN")),
                            _components.TicketViewOptions(ticketFilter),
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

func Ticket(mux *http.ServeMux, db *gorm.DB) {
    mux.HandleFunc("GET /app/ticket/{id}", func(w http.ResponseWriter, r *http.Request) {
        _middleware.MiddlewareChain(w, r,
            func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
                var ticket _model.Ticket
				id := r.PathValue("id")
				if err := db.First(&ticket, id).Error; err != nil {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
                b := NewViewBuilder("Repairs Log - Tickets", []string{
                    _components.Banner(true, _components.AppNavMenu(r.URL.Path)),
                    _components.Root(
                        _components.CenterContentWrapper(
                           _components.TicketDetails(ticket, r.URL.Query().Get("err")),
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



