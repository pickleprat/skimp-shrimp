package _partial

import (
	"net/http"

	"cfasuite/internal/_middleware"
	"cfasuite/internal/_model"

	"gorm.io/gorm"
)

func TicketList(mux *http.ServeMux, db *gorm.DB) {
    mux.HandleFunc("GET /partial/ticket/list", func(w http.ResponseWriter, r *http.Request) {
        ctx := map[string]interface{}{}
        _middleware.MiddlewareChain(ctx, w, r, _middleware.Init, _middleware.Auth,
            func(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
                var tickets []_model.Ticket
                if err := db.Find(&tickets).Error; err != nil {
                    // Handle error
                    http.Error(w, err.Error(), http.StatusInternalServerError)
                    return
                }
                w.Write()
            },
            _middleware.Log,
        )
    })
}