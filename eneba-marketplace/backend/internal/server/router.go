package server

import (
	"database/sql"
	"eneba-marketplace/backend/internal/httpx"
	"eneba-marketplace/backend/internal/middleware"
	"eneba-marketplace/backend/internal/modules/auth"
	"eneba-marketplace/backend/internal/modules/cart"
	"eneba-marketplace/backend/internal/modules/checkout"
	"eneba-marketplace/backend/internal/modules/homepage"
	"eneba-marketplace/backend/internal/modules/offers"
	"eneba-marketplace/backend/internal/modules/orders"
	"eneba-marketplace/backend/internal/modules/products"
	"net/http"
)

func Router(db *sql.DB, origin string) http.Handler {
	m := http.NewServeMux()
	p := products.Repository{DB: db}
	a := auth.Service{DB: db}
	p.Register(m)
	a.Register(m)
	offers.Register(m, db)
	homepage.Register(m, db, p)
	cart.Register(m, db)
	checkout.Service{DB: db}.Register(m, a)
	orders.Register(m, db, a)
	m.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		if db.PingContext(r.Context()) != nil {
			httpx.Error(w, 503, "Database unavailable")
			return
		}
		httpx.JSON(w, 200, map[string]string{"status": "ok"})
	})
	return middleware.CORS(origin, m)
}
