package _partial

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"cfasuite/internal/_components"
	"cfasuite/internal/_middleware"
	"cfasuite/internal/_model"

	"gorm.io/gorm"
)

func EquipmentSelectionList(mux *http.ServeMux, db *gorm.DB) {
    mux.HandleFunc("GET /partial/manufacturer/{id}/equipmentSelectionList", func(w http.ResponseWriter, r *http.Request) {
        _middleware.MiddlewareChain(w, r,
            func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
                id := r.PathValue("id")
                manufacturer := _model.Manufacturer{}
                db.First(&manufacturer, id)
                equipment := []_model.Equipment{}
				db.Where("manufacturer_id = ? AND archived = ?", id, false).Find(&equipment)
                equipmentOptions := ""
                for _, e := range equipment {
                    // Convert photo from bytes to base64
                    photoBase64 := base64.StdEncoding.EncodeToString(e.Photo)
                    // Embed photo as data URI within img tag
                    photoTag := fmt.Sprintf(`<img src="data:image/jpeg;base64,%s" alt="%s" class='w-full'>`, photoBase64, e.Nickname)
                    option := fmt.Sprintf(`
                        <div value='%d' class='equipment-option border w-[100px] hover:border-white border-darkgray cursor-pointer bg-black p-2 rounded text-sm flex items-center flex-col'>
                            %s%s
                        </div>`, e.ID, photoTag, e.Nickname)
                    equipmentOptions += option
                }
				component := ""
				if len(equipmentOptions) == 0 {
					component = fmt.Sprintf(`
						<div id='equipment-selection-list' class="flex flex-col gap-6 p-6">
							<div class='flex flex-row justify-between items-center'>
								<h2>%s has no equipment</h2>
							</div>
						</div>
					`, manufacturer.Name)
				} else {
					component = fmt.Sprintf(`
						<div id='equipment-selection-list' class="flex flex-col gap-6 p-6">
							<div class='flex flex-row justify-between items-center'>
								<h2>Assign Equipment</h2>
							</div>
							<div class='flex flex-wrap gap-4'>
								%s
							</div>
						</div>
						<script>
							document.querySelectorAll('.equipment-option').forEach(option => {
								option.addEventListener('click', (e) => {
									document.getElementById('equipment-id-input').value = option.getAttribute('value')
									document.querySelectorAll('.equipment-option').forEach(option2 => {
										option2.classList.remove('bg-white')
										option2.classList.remove('text-black')
										option2.classList.add('bg-black')
										option2.classList.add('text-white')
									})
									option.classList.remove('bg-black')
									option.classList.remove('text-white')
									option.classList.add('bg-white')
									option.classList.add('text-black')
									document.getElementById('equipment-selection-hidden-submit').click();
								})
							})
						</script>
					`, equipmentOptions)
				}
                w.Write([]byte(component))
            },
            _middleware.Init, _middleware.ParseForm, _middleware.ComponentAuth,
        )
    })
}

func EquipmentByManufacturer(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /partial/manufacturer/{id}/equipment", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				archived := r.URL.Query().Get("archived")
				manufacturer := _model.Manufacturer{}
				db.First(&manufacturer, id)
				equipment := []_model.Equipment{}
				if archived == "true" {
					db.Where("manufacturer_id = ? AND archived = ?", id, true).Find(&equipment)
					w.Write([]byte(_components.EquipmentList("Archived Equipment", equipment)))
					return
				}
				db.Where("manufacturer_id = ? AND archived = ?", id, false).Find(&equipment)
				component := _components.EquipmentList("Registered Equipment", equipment)
				w.Write([]byte(component))
			},
			_middleware.Init, _middleware.ParseForm, _middleware.ComponentAuth,
		)
	})
}

func ManufactuerList(mux *http.ServeMux, db *gorm.DB) {
    mux.HandleFunc("GET /partial/manufacturerList", func(w http.ResponseWriter, r *http.Request) {
        _middleware.MiddlewareChain(w, r,
            func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
                search := strings.ToLower(r.Form.Get("search"))
                manufacturers := []_model.Manufacturer{}
                if search == "" {
                    db.Find(&manufacturers)
                } else {
                    db.Where("LOWER(name) LIKE ?", "%"+search+"%").Find(&manufacturers)
                }
                component := _components.ManufacturerList(manufacturers)
                w.Write([]byte(component))
            },
            _middleware.Init, _middleware.ParseForm, _middleware.ComponentAuth,
        )
    })
}


func TicketList(mux *http.ServeMux, db *gorm.DB) {
    mux.HandleFunc("GET /partial/ticketList", func(w http.ResponseWriter, r *http.Request) {
        _middleware.MiddlewareChain(w, r,
            func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
                priority := r.Form.Get("priority")
                publicFilter := r.Form.Get("public")
                filteredTickets := []_model.Ticket{}
                query := db.Where("status != ?", _model.TicketStatusComplete)

                if publicFilter == "public" {
                    query = query.Where("owner <> '' AND notes <> '' AND equipment_id IS NOT NULL")
                } else if publicFilter == "private" {
                    query = query.Where("owner = '' OR notes = '' OR equipment_id IS NULL")
                }

                // Apply priority filter
                if priority != "" && priority != "all" {
                    query = query.Where("priority = ?", priority)
                }

                // Order by priority first
                query = query.Order(`
                    CASE 
                        WHEN priority = 'urgent' THEN 1
                        WHEN priority = 'medium' THEN 2
                        WHEN priority = 'low' THEN 3
                    END
                `)

                // Then order by creation date
                query = query.Order("created_at DESC")

                // Find filtered tickets
                query.Find(&filteredTickets)

                // Render component
                component := _components.TicketList(filteredTickets)
                w.Write([]byte(component))
            },
            _middleware.Init, _middleware.ParseForm, _middleware.ComponentAuth,
        )
    })
}

func PublicTicketList(mux *http.ServeMux, db *gorm.DB) {
    mux.HandleFunc("GET /partial/publicTicketList", func(w http.ResponseWriter, r *http.Request) {
        // Ensure middleware and other necessary functions are called properly.
        _middleware.MiddlewareChain(w, r, func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
            tickets := []_model.Ticket{}
            priority := r.URL.Query().Get("priority") // Correctly retrieve the priority query parameter

            // Base query to ensure all fields are filled out and status is not "complete".
            query := db.Where("creator <> '' AND item <> '' AND problem <> '' AND location <> '' AND notes <> '' AND owner <> '' AND equipment_id <> ''").
                Where("status <> ?", _model.TicketStatusComplete)

            // Modify query based on priority.
            if priority != "" && priority != "all" {
                query = query.Where("priority = ?", priority)
            }

            // Add order by clause to sort by priority.
            // This approach uses a CASE statement in raw SQL to sort priorities correctly.
            query = query.Order(`
                CASE 
                WHEN priority = 'urgent' THEN 1
                WHEN priority = 'medium' THEN 2
                WHEN priority = 'low' THEN 3
                END
            `)
			query = query.Order("created_at DESC")
            query.Find(&tickets)
            component := _components.PublicTicketList(tickets)
            w.Write([]byte(component))
        }, _middleware.Init, _middleware.ParseForm)
    })
}

func ResetEquipmentLink(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /partial/resetEquipmentLink/{ticketID}", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				ticketID := r.PathValue("ticketID")
				var ticket _model.Ticket
				var equipment _model.Equipment
				var manufacturer _model.Manufacturer
				db.First(&ticket, ticketID)
				db.First(&equipment, ticket.EquipmentID)
				db.First(&manufacturer, equipment.ManufacturerID)
				component := _components.EquipmentDetails(equipment, manufacturer)
				w.Write([]byte(component))
			},
			_middleware.Init, _middleware.ParseForm, _middleware.ComponentAuth,
		)
	})
}

func AuthWarning(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /partial/authWarning", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`
					<p>you are not authorized to access this component</p>
				`))
			},
			_middleware.Init, _middleware.ParseForm,
		)
	})
}


func CompletedTicketList(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /partial/completeTicketList/{equipmentID}", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				equipmentID := r.PathValue("equipmentID")
				tickets := []_model.Ticket{}
				db.Where("status = ? AND equipment_id = ?", _model.TicketStatusComplete, equipmentID).Find(&tickets)
				component := _components.TicketListComplete(tickets)
				w.Write([]byte(component))
			},
			_middleware.Init, _middleware.ParseForm, _middleware.ComponentAuth,
		)
	})
}


