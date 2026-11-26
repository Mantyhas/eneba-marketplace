package checkout

import (
	"eneba-marketplace/backend/internal/httpx"
	"eneba-marketplace/backend/internal/modules/auth"
	"eneba-marketplace/backend/internal/modules/cart"
	"errors"
	"net/http"
)

func (s Service) Register(m *http.ServeMux, a auth.Service) {
	m.Handle("POST /api/checkout", a.Require(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Items         []cart.Item `json:"items"`
			ExpectedTotal int         `json:"expectedTotal"`
			Payment       string      `json:"payment"`
		}
		if !httpx.Decode(w, r, &req) {
			return
		}
		if e := cart.Validate(req.Items); e != nil {
			httpx.Error(w, 400, e.Error())
			return
		}
		if req.Payment != "demo" || req.ExpectedTotal <= 0 {
			httpx.Error(w, 400, "Choose simulated payment and a valid quote")
			return
		}
		id, e := s.Purchase(r.Context(), auth.Current(r).ID, req.Items, req.ExpectedTotal)
		if errors.Is(e, ErrStock) || errors.Is(e, ErrPrice) {
			httpx.Error(w, 409, e.Error())
			return
		}
		if e != nil {
			httpx.Error(w, 500, "Purchase failed; no payment was made")
			return
		}
		httpx.JSON(w, 201, map[string]int64{"orderId": id})
	})))
}
