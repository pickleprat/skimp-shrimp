package _partial

import (
	"encoding/base64"
	"fmt"
	"net/http"

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
                    photoTag := fmt.Sprintf(`<img src="data:image/jpeg;base64,%s" alt="%s" width="50" height="50">`, photoBase64, e.Nickname)
                    option := fmt.Sprintf(`
                        <div value='%d' class='equipment-option border hover:border-white border-darkgray cursor-pointer bg-black p-2 rounded text-sm'>
                            %s %s
                        </div>`, e.ID, photoTag, e.Nickname)
                    equipmentOptions += option
                }
                component := fmt.Sprintf(`
                    <div id='equipment-selection-list' class="flex flex-col gap-6">
                        <h2>Assign Equipment</h2>
                        <div class='flex flex-wrap gap-4'>
							<a href='/app/manufacturer/%d?submitRedirect=' class='p-2 hover:border-white border border-darkgray rounded-full h-fit'><img class='h-[25px] w-[25px]' src='/static/svg/plus-dark.svg' /></a>
                            %s
                        </div>
                    </div>
                    <script>
                        document.querySelectorAll('.equipment-option').forEach(option => {
                            option.addEventListener('click', (e) => {
                                document.getElementById('equipment-id-input').value = option.getAttribute('value')
								
							})
                        })
                    </script>
                `, manufacturer.ID, equipmentOptions)
                w.Write([]byte(component))
            },
            _middleware.Init, _middleware.ParseForm, _middleware.Auth,
        )
    })
}


