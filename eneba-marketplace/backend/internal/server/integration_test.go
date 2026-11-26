package server

import (
	"bytes"
	"context"
	"encoding/json"
	"eneba-marketplace/backend/internal/database"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

// Run against an initialized disposable demo database, never a production database.
func TestPurchaseFlow(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	db, e := database.Open(context.Background(), url)
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	h := Router(db, "http://localhost:5173")
	call := func(method, path, token string, body any, status int) []byte {
		t.Helper()
		b, _ := json.Marshal(body)
		r := httptest.NewRequest(method, path, bytes.NewReader(b))
		if token != "" {
			r.Header.Set("Authorization", "Bearer "+token)
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != status {
			t.Fatalf("%s %s: got %d expected %d: %s", method, path, w.Code, status, w.Body.String())
		}
		return w.Body.Bytes()
	}
	call("GET", "/api/health", "", nil, 200)
	call("GET", "/api/homepage/banners", "", nil, 200)
	call("GET", "/api/homepage/best-sellers", "", nil, 200)
	call("GET", "/api/homepage/deals", "", nil, 200)
	var results []struct{ Slug string }
	json.Unmarshal(call("GET", "/api/products?q=eldn+ring&platform=Steam", "", nil, 200), &results)
	found := false
	for _, p := range results {
		if p.Slug == "elden-ring" {
			found = true
		}
	}
	if !found {
		t.Fatal("typo search missed Elden Ring")
	}
	call("GET", "/api/products/elden-ring", "", nil, 200)
	call("GET", "/api/orders", "", nil, 401)
	call("POST", "/api/checkout", "", map[string]any{}, 401)
	var offers []struct {
		ID         int64
		Stock      int
		PriceCents int
	}
	json.Unmarshal(call("GET", "/api/products/elden-ring/offers", "", nil, 200), &offers)
	var offerID int64
	price := 0
	for _, o := range offers {
		if o.Stock > 0 {
			offerID = o.ID
			price = o.PriceCents
			break
		}
	}
	if offerID == 0 {
		t.Fatal("seed inventory exhausted")
	}
	email := "test-" + time.Now().Format("20060102150405.000000000") + "@example.invalid"
	creds := map[string]string{"email": email, "password": "test-password-123"}
	var session struct{ Token string }
	json.Unmarshal(call("POST", "/api/auth/register", "", creds, 201), &session)
	call("POST", "/api/auth/register", "", creds, 409)
	call("POST", "/api/auth/login", "", creds, 200)
	items := []map[string]any{{"offerId": offerID, "quantity": 1}}
	call("POST", "/api/cart/quote", "", map[string]any{"items": items}, 200)
	call("POST", "/api/checkout", session.Token, map[string]any{"items": items, "expectedTotal": price + 1, "payment": "demo"}, 409)
	call("POST", "/api/checkout", session.Token, map[string]any{"items": items, "expectedTotal": price, "payment": "demo"}, 201)
	var orders []struct {
		TotalCents int
		Items      []struct{ Key string }
	}
	json.Unmarshal(call("GET", "/api/orders", session.Token, nil, 200), &orders)
	if len(orders) != 1 || orders[0].TotalCents != price || len(orders[0].Items) != 1 || orders[0].Items[0].Key == "" {
		t.Fatal("unexpected order or missing key")
	}
	call("POST", "/api/auth/logout", session.Token, map[string]any{}, 200)
	call("GET", "/api/orders", session.Token, nil, 401)
}
