package offers

import (
	"database/sql"
	"eneba-marketplace/backend/internal/httpx"
	"eneba-marketplace/backend/internal/modules/sellers"
	"net/http"
)

type Offer struct {
	ID         int64          `json:"id"`
	PriceCents int            `json:"priceCents"`
	Seller     sellers.Seller `json:"seller"`
	Stock      int            `json:"stock"`
}

func Register(m *http.ServeMux, db *sql.DB) {
	m.HandleFunc("GET /api/products/{slug}/offers", func(w http.ResponseWriter, r *http.Request) {
		rows, e := db.QueryContext(r.Context(), `SELECT o.id,o.price_cents,s.name,s.rating,(SELECT count(*) FROM digital_keys k WHERE k.offer_id=o.id AND k.sold_at IS NULL) FROM offers o JOIN sellers s ON s.id=o.seller_id JOIN products p ON p.id=o.product_id WHERE p.slug=$1 ORDER BY o.price_cents`, r.PathValue("slug"))
		if e != nil {
			httpx.Error(w, 500, "Could not load offers")
			return
		}
		defer rows.Close()
		out := []Offer{}
		for rows.Next() {
			var a Offer
			if rows.Scan(&a.ID, &a.PriceCents, &a.Seller.Name, &a.Seller.Rating, &a.Stock) != nil {
				httpx.Error(w, 500, "Could not read offers")
				return
			}
			out = append(out, a)
		}
		if rows.Err() != nil {
			httpx.Error(w, 500, "Could not read offers")
			return
		}
		httpx.JSON(w, 200, out)
	})
}
