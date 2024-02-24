package _partial

import (
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
					option := fmt.Sprintf(`
						<div value='%d' class='equipment-option border hover:bg-white hover:text-black border-white cursor-pointer bg-black p-2 rounded text-sm' value='%d'>
							%s
						</div>`, e.ID, e.ID, e.Nickname)
					equipmentOptions += option
				}
				component := fmt.Sprintf(`
					<div id='equipment-selection-list' class="flex flex-col gap-6">
						<h2>Assign Equipment</h2>
						<div class='flex flex-wrap gap-4'>
							%s
						</div>
					</div>
					<script>
						document.querySelectorAll('.equipment-option').forEach(option => {
							option.addEventListener('click', () => {
								document.getElementById('equipment-id-input').value = option.getAttribute('value')
							})
						})
					</script>
				`, equipmentOptions)
				w.Write([]byte(component))
			},
			_middleware.Init, _middleware.ParseForm, _middleware.Auth,
		)
	})
}
