package homepage

import (
    "database/sql"
    "eneba-marketplace/backend/internal/httpx"
    "eneba-marketplace/backend/internal/modules/products"
    "net/http"
    "log"
)

type Banner struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Image    string `json:"image"`
	CTA      string `json:"cta"`
}

func Register(m *http.ServeMux, db *sql.DB, p products.Repository) {
	m.HandleFunc("GET /api/homepage/banners", func(w http.ResponseWriter, r *http.Request) {
		rows, e := db.QueryContext(r.Context(), `SELECT title,subtitle,image,cta FROM banners ORDER BY position`)
		if e != nil {
    		log.Printf("Banner database error: %v", e)
   			httpx.Error(w, 500, "Could not load banners")
    		return
		
		}
		defer rows.Close()
		out := []Banner{}
		for rows.Next() {
			var b Banner
			if rows.Scan(&b.Title, &b.Subtitle, &b.Image, &b.CTA) != nil {
				httpx.Error(w, 500, "Could not read banners")
				return
			}
			out = append(out, b)
		}
		if rows.Err() != nil {
			httpx.Error(w, 500, "Could not read banners")
			return
		}
		httpx.JSON(w, 200, out)
	})
	for _, section := range []string{"best-sellers", "deals"} {
		s := section
		m.HandleFunc("GET /api/homepage/"+s, func(w http.ResponseWriter, r *http.Request) {
			sort := ""
			if s == "deals" {
				sort = "deals"
			}
			v, e := p.List(r.Context(), products.Filter{Sort: sort, Limit: 5})
			if e != nil {
				httpx.Error(w, 500, "Could not load section")
				return
			}
			httpx.JSON(w, 200, v)
		})
	}
}
