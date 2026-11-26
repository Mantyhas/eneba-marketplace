package cart

import (
	"database/sql"
	"eneba-marketplace/backend/internal/httpx"
	"errors"
	"net/http"
)

type Item struct {
	OfferID  int64 `json:"offerId"`
	Quantity int   `json:"quantity"`
}
type Request struct {
	Items []Item `json:"items"`
}

func Validate(items []Item) error {
	if len(items) == 0 || len(items) > 30 {
		return errors.New("Cart must contain 1–30 different offers")
	}
	seen := map[int64]bool{}
	for _, i := range items {
		if i.OfferID <= 0 || i.Quantity < 1 || i.Quantity > 10 || seen[i.OfferID] {
			return errors.New("Invalid or duplicate cart item")
		}
		seen[i.OfferID] = true
	}
	return nil
}

type Line struct {
	OfferID    int64  `json:"offerId"`
	Quantity   int    `json:"quantity"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	Image      string `json:"image"`
	Seller     string `json:"seller"`
	PriceCents int    `json:"priceCents"`
	Stock      int    `json:"stock"`
}
type Quote struct {
	Lines      []Line `json:"lines"`
	TotalCents int    `json:"totalCents"`
}

func Register(m *http.ServeMux, db *sql.DB) {
	m.HandleFunc("POST /api/cart/quote", func(w http.ResponseWriter, r *http.Request) {
		var req Request
		if !httpx.Decode(w, r, &req) {
			return
		}
		if len(req.Items) == 0 {
			httpx.JSON(w, 200, Quote{Lines: []Line{}})
			return
		}
		if e := Validate(req.Items); e != nil {
			httpx.Error(w, 400, e.Error())
			return
		}
		q := Quote{Lines: []Line{}}
		for _, i := range req.Items {
			l := Line{OfferID: i.OfferID, Quantity: i.Quantity}
			e := db.QueryRowContext(r.Context(), `SELECT p.name,p.slug,p.image,s.name,o.price_cents,(SELECT count(*) FROM digital_keys k WHERE k.offer_id=o.id AND k.sold_at IS NULL) FROM offers o JOIN products p ON p.id=o.product_id JOIN sellers s ON s.id=o.seller_id WHERE o.id=$1`, i.OfferID).Scan(&l.Name, &l.Slug, &l.Image, &l.Seller, &l.PriceCents, &l.Stock)
			if errors.Is(e, sql.ErrNoRows) {
				httpx.Error(w, 404, "Offer no longer exists")
				return
			}
			if e != nil {
				httpx.Error(w, 500, "Could not quote cart")
				return
			}
			q.Lines = append(q.Lines, l)
			q.TotalCents += l.PriceCents * l.Quantity
		}
		httpx.JSON(w, 200, q)
	})
}
