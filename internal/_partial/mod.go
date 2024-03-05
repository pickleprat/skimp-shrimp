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
                db.Where("manufacturer_id = ?", id).Find(&equipment)
                equipmentOptions := ""
                for _, e := range equipment {
                    // Convert photo from bytes to base64
                    photoBase64 := base64.StdEncoding.EncodeToString(e.Photo)
                    // Embed photo as data URI within img tag
                    photoTag := fmt.Sprintf(`<img src="data:image/jpeg;base64,%s" alt="%s" class='w-full'>`, photoBase64, e.Nickname)
                    option := fmt.Sprintf(`
                        <div value='%d' class='equipment-option border hover:border-white border-darkgray cursor-pointer bg-black p-2 rounded text-sm flex items-center flex-col'>
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
							<div class='grid grid-cols-3 md:grid-cols-4 gap-4'>
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
            _middleware.Init, _middleware.ParseForm, _middleware.Auth,
        )
    })
}

func EquipmentByManufacturer(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /partial/manufacturer/{id}/equipment", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				id := r.PathValue("id")
				manufacturer := _model.Manufacturer{}
				db.First(&manufacturer, id)
				equipment := []_model.Equipment{}
				db.Where("manufacturer_id = ?", id).Find(&equipment)
				component := _components.EquipmentList(equipment)
				w.Write([]byte(component))
			},
			_middleware.Init, _middleware.ParseForm, _middleware.Auth,
		)
	})
}

func Div(mux *http.ServeMux) {
	mux.HandleFunc("GET /partial/div", func(w http.ResponseWriter, r *http.Request) {
		_middleware.MiddlewareChain(w, r,
			func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				divID := r.URL.Query().Get("id")
				w.Write([]byte(fmt.Sprintf(`<div id='%s'></div>`, divID)))
			},
			_middleware.Init, _middleware.ParseForm,
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
            _middleware.Init, _middleware.ParseForm, _middleware.Auth,
        )
    })
}


func TicketList(mux *http.ServeMux, db *gorm.DB) {
    mux.HandleFunc("GET /partial/ticketList", func(w http.ResponseWriter, r *http.Request) {
        _middleware.MiddlewareChain(w, r,
            func(customContext *_middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
                priority := r.Form.Get("priority")
                status := r.Form.Get("status")
                search := strings.ToLower(r.Form.Get("search"))
                tickets := []_model.Ticket{}
                query := db
                if priority != "all" {
                    query = query.Where("priority = ?", priority)
                }
                if status != "all" {
                    query = query.Where("status = ?", status)
                }
                if search != "" {
                    query = query.Where("LOWER(creator) LIKE ? OR LOWER(item) LIKE ? OR LOWER(problem) LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
                }
                query.Find(&tickets)
                component := _components.TicketList(tickets)
                w.Write([]byte(component))
            },
            _middleware.Init, _middleware.ParseForm, _middleware.Auth,
        )
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
			_middleware.Init, _middleware.ParseForm, _middleware.Auth,
		)
	})
}



