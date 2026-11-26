package products

import (
	"database/sql"
	"eneba-marketplace/backend/internal/httpx"
	"errors"
	"net/http"
	"strconv"
)

func (r Repository) Register(m *http.ServeMux) {
	m.HandleFunc("GET /api/products", r.list)
	m.HandleFunc("GET /api/products/{slug}", r.get)
}
func (r Repository) list(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	max, e := strconv.Atoi(q.Get("max"))
	if q.Get("max") != "" && (e != nil || max < 0) {
		httpx.Error(w, 400, "Invalid maximum price")
		return
	}
	v, e := r.List(req.Context(), Filter{Query: q.Get("q"), Platform: q.Get("platform"), Region: q.Get("region"), Sort: q.Get("sort"), Max: max})
	if e != nil {
		httpx.Error(w, 500, "Could not load products")
		return
	}
	httpx.JSON(w, 200, v)
}
func (r Repository) get(w http.ResponseWriter, req *http.Request) {
	p, e := r.Get(req.Context(), req.PathValue("slug"))
	if errors.Is(e, sql.ErrNoRows) {
		httpx.Error(w, 404, "Product not found")
		return
	}
	if e != nil {
		httpx.Error(w, 500, "Could not load product")
		return
	}
	httpx.JSON(w, 200, p)
}
