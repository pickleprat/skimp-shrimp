package _view

import (
	"cfasuite/internal/_components"
	"cfasuite/internal/_middleware"
	"cfasuite/internal/_model"
	"cfasuite/internal/_util"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
		<body hx-boost='true' hx-history='false' class='bg-black text-white'>
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

func Login(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				sessionCookie, _ := r.Cookie("SessionToken")
				if sessionCookie != nil && sessionCookie.Value == os.Getenv("ADMIN_SESSION_TOKEN") {
					http.Redirect(w, r, "/app", http.StatusSeeOther)
					return
				}
				b := NewViewBuilder("Repairs Log - Login", []string{
					_components.MainLoader(),
					_components.Root(
						_components.CenterContentWrapper(
							_components.Banner("Repairs Log", ""),
							_components.LoginForm(r.URL.Query().Get("err"), r.URL.Query().Get("username"), r.URL.Query().Get("password")),
						),
					),
					_components.BottomSpacer(),
				})
				w.Write(b.Build())
			},
			_middleware.Init,
		)
	})
}

func AdminHome(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				b := NewViewBuilder("Repairs Log - Home", []string{
					_components.MainLoader(),
					_components.Root(
						_components.CenterContentWrapper(
							_components.Banner("Repairs Log", 
								_components.SimpleTopNav(
									_components.NavLink("Home", "/app", true),
									_components.NavLink("Manufacturers", "/app/manufacturer", false),
									_components.NavLink("Tickets", "/app/ticket", false),
									_util.ConditionalString(
										os.Getenv("GO_ENV") == "dev",
										_components.NavLink("Public View", _util.URLBuilder("/app/ticket/public", "publicSecurityToken", os.Getenv("PUBLIC_SECURITY_TOKEN")), false),
										"",
									),
									_components.NavLink("Logout", "/logout", false),
								),
							),
							_components.TitleAndText("Welcome Home!", "Repairs Log is an application intended to help track equipment repairs over a long period of time. To start, you can begin registering manufacturers for your business."),
							_components.TitleAndText("Manufacturers", "Everything starts with manufacturers. In business, we often need to reach out to a manufacturer to assist with a repair, to order a part, or to get advice. Now you can spend less time searching for their contact info and more time building your business."),
							_components.TitleAndText("Equipment", "Once you have manufacturers registered, you can begin registering equipment. When registering a piece of equipment, you will assign it to a manufacturer. This will give you quick access to important details like model numbers, serial numbers, and the equipment's repair history."),
							_components.TitleAndText("Tickets", "Finally, you can begin creating tickets. Tickets are the lifeblood of this application. They help you keep track of what needs to be repaired, what has been repaired, and what is in good working order. Tickets can be created by your team and then assigned to a piece of equipment within the application. From there, you can update the public details of a ticket to allow all members or your team to see the status of a particular repair."),
						),
					),
					_components.BottomSpacer(),
				})
				w.Write(b.Build())
			},
			_middleware.Init, _middleware.Auth,
		)
	})
}

func CreateManufacturers(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/manufacturer", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				var manufacturers []_model.Manufacturer
				submitRedirect := r.URL.Query().Get("submitRedirect")
				db.Find(&manufacturers)
				b := NewViewBuilder("Repairs Log - Manufacturers", []string{
					_components.MainLoader(),
					_components.Root(
						_components.CenterContentWrapper(
							_components.Banner("Repairs Log", 
							_components.BreadCrumbs(
									_components.NavLink("Home", "/app", false),
									_components.NavLink("Manufacturers", "/app/manufacturer", true),
								),
							),
							_components.CreateManufacturerForm(r.URL.Query().Get("err"), r.URL.Query().Get("success"), r.URL.Query().Get("name"), r.URL.Query().Get("phone"), r.URL.Query().Get("email"), submitRedirect),
							`
							<form action='/partial/manufacturerList' method='GET' class="p-6">
								<div id='manufacturer-search' class='flex justify-between bg-black border border-darkgray p-2 rounded-lg'>
									<div class='h-6 w-6'>
										<img src="/static/svg/search-dark.svg">
									</div>
									<div class='flex flex-row justify-between items-center w-full'>
										<input name='search' id='manufacturer-search-input' 
											hx-target='#manufacturer-list' 
											hx-get='/partial/manufacturerList' 
											hx-swap='outerHTML' 
											hx-trigger='search, keyup delay:200ms changed'
											hx-indicator='#manufacturer-list-indicator' 
											type="text" 
											class="pl-4 w-full bg-inherit focus:outline-none" 
											placeholder="Search Manufacturers...">
										<div id='manufacturer-list-indicator' class='flex-indicator animate-spin rounded-full border border-darkgray border-t-white h-6 w-6 p-2'></div>
									</div>
								</div>
							</form>
							<script>
								document.body.addEventListener = function() {};
								document.body.addEventListener('click', (event) => {
									if (qs('#manufacturer-search')) {
										if (!qs('#manufacturer-search').contains(event.target)) {
											qs('#manufacturer-search').classList.remove('border-lightgray');
										}
									}
								});
								qs('#manufacturer-search').addEventListener('click', (event) => {
									event.stopPropagation();
									qs('#manufacturer-search').classList.toggle('border-lightgray');
								});
							</script>
							`,
							_components.HxGetLoader("/partial/manufacturerList"),
						),
					),
					_components.BottomSpacer(),
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
                var equipment []_model.Equipment
                
                // Check if the manufacturer exists in the database
                if err := db.First(&manufacturer, id).Error; err != nil {
                    if errors.Is(err, gorm.ErrRecordNotFound) {
                        http.Error(w, "Manufacturer not found", http.StatusBadRequest)
                        return
                    }
                    // Handle other database errors
                    http.Error(w, "Internal server error", http.StatusInternalServerError)
                    return
                }
                db.Where("manufacturer_id = ?", manufacturer.ID).Find(&equipment)
                b := NewViewBuilder("Repairs Log - " + manufacturer.Name, []string{
                    _components.MainLoader(),
                    _components.Root(
                        _components.CenterContentWrapper(
                            _components.Banner("Repairs Log", 
                                _components.BreadCrumbs(
                                    _components.NavLink("Home", "/app", false),
                                    _components.NavLink("Manufacturers", "/app/manufacturer", false),
                                    _components.NavLink("Equipment", fmt.Sprintf("/app/manufacturer/%d", manufacturer.ID), true),
                                ),
                            ),
                            _components.TopNav(
                                _components.LinkButton(fmt.Sprintf("/app/manufacturer/%d", manufacturer.ID), "Equipment", true),
                                _components.LinkButton(fmt.Sprintf("/app/manufacturer/%d/archive", manufacturer.ID), "Archived Equipment", false),
                                _components.LinkButton(fmt.Sprintf("/app/manufacturer/%d/update", manufacturer.ID), "Update", false),
                                _components.LinkButton(fmt.Sprintf("/app/manufacturer/%d/delete", manufacturer.ID), "Delete", false),
                            ),
                            _components.ManufacturerDetails(manufacturer),
                            _components.CreateEquipmentForm(
                                manufacturer, 
                                r.URL.Query().Get("err"), 
                                r.URL.Query().Get("success"),
                                r.URL.Query().Get("submitRedirect"),
                                r.URL.Query().Get("nickname"),
                                r.URL.Query().Get("serialNumber"),
                                r.URL.Query().Get("modelNumber"),
                            ),
                            _components.HxGetLoader(
                                fmt.Sprintf("/partial/manufacturer/%d/equipment", manufacturer.ID),
                            ),
                        ),
                    ),
                    _components.BottomSpacer(),
                })
                w.Write(b.Build())
            },
            _middleware.Init, _middleware.Auth,
        )
    })
}


func EquipmentArchive(mux *http.ServeMux, db *gorm.DB) {
    mux.HandleFunc("GET /app/manufacturer/{id}/archive", func(w http.ResponseWriter, r *http.Request) {
        _middleware.MiddlewareChain(w, r,
            func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
                id := r.PathValue("id")
                var manufacturer _model.Manufacturer
                var equipment []_model.Equipment
                if err := db.First(&manufacturer, id).Error; err != nil {
                    http.Error(w, "Manufacturer not found", http.StatusBadRequest)
                    return
                }
                db.Where("manufacturer_id = ?", manufacturer.ID).Find(&equipment)
                b := NewViewBuilder("Repairs Log - " + manufacturer.Name, []string{
                    _components.MainLoader(),
                    _components.Root(
                        _components.CenterContentWrapper(
                            _components.Banner("Repairs Log", 
                                _components.BreadCrumbs(
                                    _components.NavLink("Home", "/app", false),
                                    _components.NavLink("Manufacturers", "/app/manufacturer", false),
                                    _components.NavLink("Equipment Archive", fmt.Sprintf("/app/manufacturer/%d/archive", manufacturer.ID), true),
                                ),
                            ),
                            _components.TopNav(
                                _components.LinkButton(fmt.Sprintf("/app/manufacturer/%d", manufacturer.ID), "Equipment", false),
                                _components.LinkButton(fmt.Sprintf("/app/manufacturer/%d/archive", manufacturer.ID), "Archived Equipment", true),
                                _components.LinkButton(fmt.Sprintf("/app/manufacturer/%d/update", manufacturer.ID), "Update", false),
                                _components.LinkButton(fmt.Sprintf("/app/manufacturer/%d/delete", manufacturer.ID), "Delete", false),
                            ),
                            _components.ManufacturerDetails(manufacturer),
                            _components.HxGetLoader(
                                fmt.Sprintf("/partial/manufacturer/%d/equipment?archived=true", manufacturer.ID),
                            ),
                        ),
                    ),
                    _components.BottomSpacer(),
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
                if err := db.First(&manufacturer, id).Error; err != nil {
                    http.Error(w, "Manufacturer not found", http.StatusBadRequest)
                    return
                }
                b := NewViewBuilder("Repairs Log - " + manufacturer.Name, []string{
                    _components.MainLoader(),
                    _components.Root(
                        _components.CenterContentWrapper(
                            _components.Banner("Repairs Log",
                                _components.BreadCrumbs(
                                    _components.NavLink("Home", "/app", false),
                                    _components.NavLink("Manufacturers", "/app/manufacturer", false),
                                    _components.NavLink("Delete", fmt.Sprintf("/app/manufacturer/%d/delete", manufacturer.ID), true),
                                ),
                            ),
                            _components.TopNav(
                                _components.LinkButton(fmt.Sprintf("/app/manufacturer/%d", manufacturer.ID), "Equipment", false),
                                _components.LinkButton(fmt.Sprintf("/app/manufacturer/%d/archive", manufacturer.ID), "Archived Equipment", false),
                                _components.LinkButton(fmt.Sprintf("/app/manufacturer/%d/update", manufacturer.ID), "Update", false),
                                _components.LinkButton(fmt.Sprintf("/app/manufacturer/%d/delete", manufacturer.ID), "Delete", true),
                            ),
                            _components.ManufacturerDetails(manufacturer),
                            _components.DeleteManufacturerForm(manufacturer, r.URL.Query().Get("err")),
                        ),
                    ),
                    _components.BottomSpacer(),
                })
                w.Write(b.Build())
            },
            _middleware.Init, _middleware.Auth,
        )
    })
}


func UpdateManufacturer(mux *http.ServeMux, db *gorm.DB) {
    mux.HandleFunc("GET /app/manufacturer/{id}/update", func(w http.ResponseWriter, r *http.Request) {
        _middleware.MiddlewareChain(w, r,
            func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
                id := r.PathValue("id")
                var manufacturer _model.Manufacturer
                if err := db.First(&manufacturer, id).Error; err != nil {
                    http.Error(w, "Manufacturer not found", http.StatusBadRequest)
                    return
                }
                b := NewViewBuilder("Repairs Log - " + manufacturer.Name, []string{
                    _components.MainLoader(),
                    _components.Root(
                        _components.CenterContentWrapper(
                            _components.Banner("Repairs Log", 
                                _components.BreadCrumbs(
                                    _components.NavLink("Home", "/app", false),
                                    _components.NavLink("Manufacturers", "/app/manufacturer", false),
                                    _components.NavLink("Update", fmt.Sprintf("/app/manufacturer/%d/update", manufacturer.ID), true),
                                ),
                            ),
                            _components.TopNav(
                                _components.LinkButton(fmt.Sprintf("/app/manufacturer/%d", manufacturer.ID), "Equipment", false),
                                _components.LinkButton(fmt.Sprintf("/app/manufacturer/%d/archive", manufacturer.ID), "Archived Equipment", false),
                                _components.LinkButton(fmt.Sprintf("/app/manufacturer/%d/update", manufacturer.ID), "Update", true),
                                _components.LinkButton(fmt.Sprintf("/app/manufacturer/%d/delete", manufacturer.ID), "Delete", false),
                            ),
                            _components.ManufacturerDetails(manufacturer),
                            _components.UpdateManufacturerForm(manufacturer, r.URL.Query().Get("err"), r.URL.Query().Get("success")),
                        ),
                    ),
                    _components.BottomSpacer(),
                })
                w.Write(b.Build())
            },
            _middleware.Init, _middleware.Auth,
        )
    })
}


func UpdateEquipment(mux *http.ServeMux, db *gorm.DB) {
    mux.HandleFunc("GET /app/equipment/{equipmentID}", func(w http.ResponseWriter, r *http.Request) {
        _middleware.MiddlewareChain(w, r,
            func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
                equipmentID := r.PathValue("equipmentID")
                var equipment _model.Equipment
                var manufacturer _model.Manufacturer
                if err := db.First(&equipment, equipmentID).Error; err != nil {
                    http.Error(w, "Equipment not found", http.StatusBadRequest)
                    return
                }
                if err := db.First(&manufacturer, equipment.ManufacturerID).Error; err != nil {
                    http.Error(w, "Manufacturer not found", http.StatusInternalServerError)
                    return
                }
                b := NewViewBuilder(fmt.Sprintf("Repairs Log - %s", equipment.Nickname), []string{
                    _components.MainLoader(),
                    _components.Root(
                        _components.CenterContentWrapper(
                            _components.Banner("Repairs Log",
                                _components.BreadCrumbs(
                                    _components.NavLink("Home", "/app", false),
                                    _components.NavLink("Manufacturers", "/app/manufacturer", false),
                                    _components.NavLink("Equipment", fmt.Sprintf("/app/manufacturer/%d", manufacturer.ID), false),
                                    _components.NavLink("Update", fmt.Sprintf("/app/equipment/%d", equipment.ID), true),
                                ),
                            ),
                            _components.TopNav(
                                _components.LinkButton(fmt.Sprintf("/app/equipment/%s", equipmentID), "Update", true),
                                _components.LinkButton(fmt.Sprintf("/app/equipment/%s/delete", equipmentID), "Delete", false),
                            ),
                            _components.EquipmentDetails(equipment, manufacturer),
                            _components.UpdateEquipmentForm(equipment, r.URL.Query().Get("err"), r.URL.Query().Get("success")),
							_components.HxGetLoader(fmt.Sprintf("/partial/completeTicketList/%s", equipmentID)),
                        ),
                    ),
                    _components.BottomSpacer(),
                })
                w.Write(b.Build())
            },
            _middleware.Init, _middleware.Auth,
        )
    })
}


func DeleteEquipment(mux *http.ServeMux, db *gorm.DB) {
    mux.HandleFunc("GET /app/equipment/{equipmentID}/delete", func(w http.ResponseWriter, r *http.Request) {
        _middleware.MiddlewareChain(w, r,
            func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
                equipmentID := r.PathValue("equipmentID")
                var equipment _model.Equipment
                var manufacturer _model.Manufacturer
                if err := db.First(&equipment, equipmentID).Error; err != nil {
                    http.Error(w, "Equipment not found", http.StatusBadRequest)
                    return
                }
                if err := db.First(&manufacturer, equipment.ManufacturerID).Error; err != nil {
                    http.Error(w, "Manufacturer not found", http.StatusInternalServerError)
                    return
                }
                b := NewViewBuilder(fmt.Sprintf("Repairs Log - %s", equipment.Nickname), []string{
                    _components.MainLoader(),
                    _components.Root(
                        _components.CenterContentWrapper(
                            _components.Banner("Repairs Log", 
                                _components.BreadCrumbs(
                                    _components.NavLink("Home", "/app", false),
                                    _components.NavLink("Manufacturers", "/app/manufacturer", false),
                                    _components.NavLink("Equipment", fmt.Sprintf("/app/manufacturer/%d", manufacturer.ID), false),
                                    _components.NavLink("Delete", fmt.Sprintf("/app/equipment/%d/delete", equipment.ID), true),
                                ),
                            ),
                            _components.TopNav(
                                _components.LinkButton(fmt.Sprintf("/app/equipment/%s", equipmentID), "Update", false),
                                _components.LinkButton(fmt.Sprintf("/app/equipment/%s/delete", equipmentID), "Delete", true),
                            ),
                            _components.EquipmentDetails(equipment, manufacturer),
                            _components.DeleteEquipmentForm(equipment, r.URL.Query().Get("err"), r.URL.Query().Get("success")),
                        ),
                    ),
                    _components.BottomSpacer(),
                })
                w.Write(b.Build())
            },
            _middleware.Init, _middleware.Auth,
        )
    })
}



func PublicCreateTicket(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/ticket/public", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				token := r.URL.Query().Get("publicSecurityToken")
				if token != os.Getenv("PUBLIC_SECURITY_TOKEN") {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
				b := NewViewBuilder("Repairs Log - Public Tickets", []string{
					_components.MainLoader(),
					_components.Root(
						_components.CenterContentWrapper(
							_components.Banner("Repairs Log", 
								_components.SimpleTopNav(
									_components.NavLink("Create Tickets", _util.URLBuilder("/app/ticket/public", "publicSecurityToken", token), true),
									_components.NavLink("View Tickets", _util.URLBuilder("/app/ticket/public/view", "publicSecurityToken", token), false),
								),
							),
							_components.PublicCreateTicketForm(r, token, "/form/ticket/public"),
						),
					),
					_components.BottomSpacer(),
				})
				w.Write(b.Build())
			},
			_middleware.Init,
		)
	})
}

func PublicViewTickets(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/ticket/public/view", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				token := r.URL.Query().Get("publicSecurityToken")
				if token != os.Getenv("PUBLIC_SECURITY_TOKEN") {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
				b := NewViewBuilder("Repairs Log - Public Tickets", []string{
					_components.MainLoader(),
					_components.Root(
						_components.CenterContentWrapper(
							_components.Banner("Repairs Log", 
								_components.SimpleTopNav(
									_components.NavLink("Create Tickets", _util.URLBuilder("/app/ticket/public", "publicSecurityToken", token), false),
									_components.NavLink("View Tickets", _util.URLBuilder("/app/ticket/public/view", "publicSecurityToken", token), true),
								),
							),
							_components.PublicTicketSearchForm(),
							_components.HxGetLoader(_util.URLBuilder("/partial/publicTicketList", "publicSecurityToken", token)),
						),
					),
					_components.BottomSpacer(),
				})
				w.Write(b.Build())
			},
			_middleware.Init,
		)
	})

}

func AdminViewTicket(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/ticket/{id}", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				var ticket _model.Ticket
				id := r.PathValue("id")
				if err := db.First(&ticket, id).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						http.Error(w, "Ticket not found", http.StatusBadRequest)
						return
					}
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				var manufacturers []_model.Manufacturer
				if ticket.EquipmentID == nil {
					db.Find(&manufacturers)
				}
				b := NewViewBuilder("Repairs Log - Assign Equipment", []string{
					_components.MainLoader(),
					_components.Root(
						_components.CenterContentWrapper(
							_components.Banner("Repairs Log",
								_components.BreadCrumbs(
									_components.NavLink("Home", "/app", false),
									_components.NavLink("Tickets", "/app/ticket", false),
									_components.NavLink("Assign Equipment", fmt.Sprintf("/app/ticket/%d", ticket.ID), true),
								),
							),
							_components.TopNav(
								_components.LinkButton(fmt.Sprintf("/app/ticket/%d", ticket.ID), "Assign Equipment", true),
								_components.LinkButton(fmt.Sprintf("/app/ticket/%d/update", ticket.ID), "Update", false),
								_util.ConditionalString(
									ticket.Owner == "" || ticket.Notes == "" || ticket.EquipmentID == nil,
									"",
									_components.LinkButton(fmt.Sprintf("/app/ticket/%d/complete", ticket.ID), "Complete", false),
								),
								_components.LinkButton(fmt.Sprintf("/app/ticket/%d/delete", ticket.ID), "Delete", false),
							),
							_components.TicketDetails(ticket, ""),
							_util.ConditionalString(
								ticket.EquipmentID != nil,
								fmt.Sprintf(`
									<div class='flex p-6 flex-col gap-6'>
										<h2 class='text-lg'>Equipment Association</h2>
										<p class='text-xs'>This ticket is assigned to the equipment shown below. Click the following button to remove the association:</p>
										<a href='/form/ticket/%d/resetEquipment' class='w-full p-2 bg-white text-center text-black rounded'>Remove Equipment Association</a>
									</div>
								`, ticket.ID),
								"",
							),
							_util.ConditionalString(
								ticket.EquipmentID == nil,
								_components.TicketAssignmentForm(ticket, manufacturers, db),
								_components.HxGetLoader("/partial/resetEquipmentLink/"+id),
							),
						),
					),
					_components.BottomSpacer(),
				})
				w.Write(b.Build())
			},
			_middleware.Init, _middleware.Auth,
		)
	})
}


func AdminCreateTickets(mux *http.ServeMux, db *gorm.DB) {
    mux.HandleFunc("GET /app/ticket", func(w http.ResponseWriter, r *http.Request) {
        _middleware.MiddlewareChain(w, r,
            func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
                b := NewViewBuilder("Repairs Log - Tickets", []string{
					_components.MainLoader(),
                    _components.Root(
						_components.CenterContentWrapper(
							_components.Banner("Repairs Log", 
								_components.BreadCrumbs(
									_components.NavLink("Home", "/app", false),
									_components.NavLink("Tickets", "/app/ticket", true),
								),
							),
                            _components.CreateTicketForm(r, "", "/form/ticket/admin"),
							_components.TicketSearchForm(),
                            _components.HxGetLoader("/partial/ticketList?public=private&priority=all"),
                        ),
                    ),
                    _components.BottomSpacer(),
                })
                w.Write(b.Build())
            },
            _middleware.Init, _middleware.Auth,
        )
    })
}

func AdminUpdateTicket(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/ticket/{id}/update", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				var ticket _model.Ticket
				id := r.PathValue("id")
				if err := db.First(&ticket, id).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						http.Error(w, "Ticket not found", http.StatusBadRequest)
						return
					}
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				b := NewViewBuilder("Repairs Log - Update Ticket", []string{
					_components.MainLoader(),
					_components.Root(
						_components.CenterContentWrapper(
							_components.Banner("Repairs Log", 
								_components.BreadCrumbs(
									_components.NavLink("Home", "/app", false),
									_components.NavLink("Tickets", "/app/ticket", false),
									_components.NavLink("Update", fmt.Sprintf("/app/ticket/%d/update", ticket.ID), true),
								),
							),
							_components.TopNav(
								_components.LinkButton(fmt.Sprintf("/app/ticket/%d", ticket.ID), "Assign Equipment", false),
								_components.LinkButton(fmt.Sprintf("/app/ticket/%d/update", ticket.ID), "Update", true),
								_util.ConditionalString(
									ticket.Owner == "" || ticket.Notes == "" || ticket.EquipmentID == nil,
									"",
									_components.LinkButton(fmt.Sprintf("/app/ticket/%d/complete", ticket.ID), "Complete", false),
								),
								_components.LinkButton(fmt.Sprintf("/app/ticket/%d/delete", ticket.ID), "Delete", false),
							),
							_components.UpdateTicketForm(ticket, r.URL.Query().Get("err"), r.URL.Query().Get("success")),
						),
					),
					_components.BottomSpacer(),
				})
				w.Write(b.Build())
			},
			_middleware.Init, _middleware.Auth,
		)
	})
}


func AdminDeleteTicket(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/ticket/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				var ticket _model.Ticket
				id := r.PathValue("id")
				if err := db.First(&ticket, id).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						http.Error(w, "Ticket not found", http.StatusBadRequest)
						return
					}
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				b := NewViewBuilder("Repairs Log - Delete Ticket", []string{
					_components.MainLoader(),
					_components.Root(
						_components.CenterContentWrapper(
							_components.Banner("Repairs Log", 
								_components.BreadCrumbs(
									_components.NavLink("Home", "/app", false),
									_components.NavLink("Tickets", "/app/ticket", false),
									_components.NavLink("Delete", fmt.Sprintf("/app/ticket/%d/delete", ticket.ID), true),
								),
							),
							_components.TopNav(
								_components.LinkButton(fmt.Sprintf("/app/ticket/%d", ticket.ID), "Assign Equipment", false),
								_components.LinkButton(fmt.Sprintf("/app/ticket/%d/update", ticket.ID), "Update", false),
								_util.ConditionalString(
									ticket.Owner == "" || ticket.Notes == "" || ticket.EquipmentID == nil,
									"",
									_components.LinkButton(fmt.Sprintf("/app/ticket/%d/complete", ticket.ID), "Complete", false),
								),
								_components.LinkButton(fmt.Sprintf("/app/ticket/%d/delete", ticket.ID), "Delete", true),
							),
							_components.DeleteTicketForm(ticket, r.URL.Query().Get("err")),
						),
					),
					_components.BottomSpacer(),
				})
				w.Write(b.Build())
			},
			_middleware.Init, _middleware.Auth,
		)
	})
}


func AdminCompleteTicket(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /app/ticket/{id}/complete", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				var ticket _model.Ticket
				id := r.PathValue("id")
				if err := db.First(&ticket, id).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						http.Error(w, "Ticket not found", http.StatusBadRequest)
						return
					}
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				b := NewViewBuilder("Repairs Log - Complete Ticket", []string{
					_components.MainLoader(),
					_components.Root(
						_components.CenterContentWrapper(
							_components.Banner("Repairs Log", 
								_components.BreadCrumbs(
									_components.NavLink("Home", "/app", false),
									_components.NavLink("Tickets", "/app/ticket", false),
									_components.NavLink("Complete", fmt.Sprintf("/app/ticket/%d/complete", ticket.ID), true),
								),
							),
							_components.TopNav(
								_components.LinkButton(fmt.Sprintf("/app/ticket/%d", ticket.ID), "Assign Equipment", false),
								_components.LinkButton(fmt.Sprintf("/app/ticket/%d/update", ticket.ID), "Update", false),
								_components.LinkButton(fmt.Sprintf("/app/ticket/%d/complete", ticket.ID), "Complete", true),
								_components.LinkButton(fmt.Sprintf("/app/ticket/%d/delete", ticket.ID), "Delete", false),
							),
							_components.CompleteTicketForm(ticket, r.URL.Query().Get("err"), r.URL.Query().Get("success")),
						),
					),
					_components.BottomSpacer(),
				})
				w.Write(b.Build())
			},
			_middleware.Init, _middleware.Auth,
		)
	})
}

