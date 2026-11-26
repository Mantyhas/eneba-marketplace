package checkout

import (
	"context"
	"database/sql"
	"eneba-marketplace/backend/internal/modules/cart"
	"errors"
	"sort"
)

var ErrStock = errors.New("An offer is unavailable or has insufficient stock")
var ErrPrice = errors.New("Price changed; review your cart and try again")

type Service struct{ DB *sql.DB }

// Purchase rechecks prices and locks inventory inside a single transaction.
func (s Service) Purchase(ctx context.Context, user int64, items []cart.Item, expected int) (int64, error) {
	if e := cart.Validate(items); e != nil {
		return 0, e
	}
	items = append([]cart.Item(nil), items...)
	sort.Slice(items, func(i, j int) bool { return items[i].OfferID < items[j].OfferID })
	tx, e := s.DB.BeginTx(ctx, nil)
	if e != nil {
		return 0, e
	}
	defer tx.Rollback()
	var order int64
	e = tx.QueryRowContext(ctx, `INSERT INTO orders(user_id,total_cents) VALUES($1,0) RETURNING id`, user).Scan(&order)
	if e != nil {
		return 0, e
	}
	total := 0
	for _, i := range items {
		var price int
		e = tx.QueryRowContext(ctx, `SELECT price_cents FROM offers WHERE id=$1 FOR UPDATE`, i.OfferID).Scan(&price)
		if errors.Is(e, sql.ErrNoRows) {
			return 0, ErrStock
		}
		if e != nil {
			return 0, e
		}
		for n := 0; n < i.Quantity; n++ {
			var key int64
			e = tx.QueryRowContext(ctx, `SELECT id FROM digital_keys WHERE offer_id=$1 AND sold_at IS NULL ORDER BY id LIMIT 1 FOR UPDATE`, i.OfferID).Scan(&key)
			if errors.Is(e, sql.ErrNoRows) {
				return 0, ErrStock
			}
			if e != nil {
				return 0, e
			}
			if _, e = tx.ExecContext(ctx, `UPDATE digital_keys SET sold_at=now() WHERE id=$1`, key); e != nil {
				return 0, e
			}
			if _, e = tx.ExecContext(ctx, `INSERT INTO order_items(order_id,offer_id,key_id,price_cents) VALUES($1,$2,$3,$4)`, order, i.OfferID, key, price); e != nil {
				return 0, e
			}
			total += price
		}
	}
	if total != expected {
		return 0, ErrPrice
	}
	if _, e = tx.ExecContext(ctx, `UPDATE orders SET total_cents=$1 WHERE id=$2`, total, order); e != nil {
		return 0, e
	}
	return order, tx.Commit()
}
