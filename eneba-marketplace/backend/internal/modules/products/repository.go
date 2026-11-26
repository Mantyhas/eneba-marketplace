package products

import (
	"context"
	"database/sql"
	"eneba-marketplace/backend/internal/modules/search"
)

type Repository struct{ DB *sql.DB }

const base = `SELECT p.id,p.slug,p.name,p.platform,p.region,p.description,p.image,p.original_cents,COALESCE(MIN(o.price_cents) FILTER (WHERE EXISTS(SELECT 1 FROM digital_keys k WHERE k.offer_id=o.id AND k.sold_at IS NULL)),0),p.popularity FROM products p LEFT JOIN offers o ON o.product_id=p.id`

func (r Repository) List(ctx context.Context, f Filter) ([]Product, error) {
	order := "p.popularity DESC,p.id"
	switch f.Sort {
	case "price_asc":
		order = "9 ASC,p.id"
	case "price_desc":
		order = "9 DESC,p.id"
	case "name":
		order = "p.name ASC,p.id"
	case "deals":
		order = "(p.original_cents-COALESCE(MIN(o.price_cents) FILTER (WHERE EXISTS(SELECT 1 FROM digital_keys k WHERE k.offer_id=o.id AND k.sold_at IS NULL)),p.original_cents)) DESC,p.id"
	}
	if f.Limit <= 0 || f.Limit > 200 {
		f.Limit = 200
	}
	rows, e := r.DB.QueryContext(ctx, base+` WHERE ($1='' OR lower(p.name) LIKE $2 OR word_similarity(lower($1),lower(p.name)) >= 0.38) AND ($3='' OR p.platform=$3) AND ($4='' OR p.region=$4) GROUP BY p.id HAVING ($5=0 OR COALESCE(MIN(o.price_cents) FILTER (WHERE EXISTS(SELECT 1 FROM digital_keys k WHERE k.offer_id=o.id AND k.sold_at IS NULL)),0) BETWEEN 1 AND $5) ORDER BY `+order+` LIMIT $6`, f.Query, search.Pattern(f.Query), f.Platform, f.Region, f.Max, f.Limit)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []Product{}
	for rows.Next() {
		var p Product
		if e = scan(rows, &p); e != nil {
			return nil, e
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

type scanner interface{ Scan(...any) error }

func scan(s scanner, p *Product) error {
	return s.Scan(&p.ID, &p.Slug, &p.Name, &p.Platform, &p.Region, &p.Description, &p.Image, &p.OriginalCents, &p.PriceCents, &p.Popularity)
}
func (r Repository) Get(ctx context.Context, slug string) (Product, error) {
	var p Product
	e := scan(r.DB.QueryRowContext(ctx, base+` WHERE p.slug=$1 GROUP BY p.id`, slug), &p)
	return p, e
}
