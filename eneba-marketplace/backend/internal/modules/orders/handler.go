package orders

import (
	"database/sql"
	"eneba-marketplace/backend/internal/httpx"
	"eneba-marketplace/backend/internal/modules/auth"
	"net/http"
	"time"
)

type Item struct {
	Name       string `json:"name"`
	Key        string `json:"key"`
	PriceCents int    `json:"priceCents"`
}
type Order struct {
	ID         int64     `json:"id"`
	TotalCents int       `json:"totalCents"`
	CreatedAt  time.Time `json:"createdAt"`
	Items      []Item    `json:"items"`
}

func Register(m *http.ServeMux, db *sql.DB, a auth.Service) {
	m.Handle("GET /api/orders", a.Require(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		rows, e := db.QueryContext(r.Context(), `SELECT o.id,o.total_cents,o.created_at,p.name,k.code,i.price_cents FROM orders o JOIN order_items i ON i.order_id=o.id JOIN digital_keys k ON k.id=i.key_id JOIN offers f ON f.id=i.offer_id JOIN products p ON p.id=f.product_id WHERE o.user_id=$1 ORDER BY o.id DESC,i.id`, auth.Current(r).ID)
		if e != nil {
			httpx.Error(w, 500, "Could not load orders")
			return
		}
		defer rows.Close()
		out := []Order{}
		for rows.Next() {
			var o Order
			var i Item
			if rows.Scan(&o.ID, &o.TotalCents, &o.CreatedAt, &i.Name, &i.Key, &i.PriceCents) != nil {
				httpx.Error(w, 500, "Could not read orders")
				return
			}
			if len(out) == 0 || out[len(out)-1].ID != o.ID {
				out = append(out, o)
			}
			out[len(out)-1].Items = append(out[len(out)-1].Items, i)
		}
		if rows.Err() != nil {
			httpx.Error(w, 500, "Could not read orders")
			return
		}
		httpx.JSON(w, 200, out)
	})))
}
