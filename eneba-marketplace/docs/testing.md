# Test plan

## Executed in build environment
- `go test ./internal/modules/cart ./internal/modules/search ./internal/httpx ./internal/middleware`: passed (cart and search assertions; HTTP and middleware had no tests).
- Go source formatted with gofmt.
- TypeScript transpile syntax checks passed; no dependency-aware typecheck or bundle build.

## Pending local verification
Run `go mod tidy; go test ./...` in backend and `npm install; npm run build` in frontend. The integration test is opt-in and consumes one seeded demo key. Use a disposable initialized database, not production:

```powershell
$env:TEST_DATABASE_URL = 'postgres://postgres:YOUR_PASSWORD@localhost:5432/eneba_marketplace?sslmode=disable'
go test ./internal/server -run TestPurchaseFlow -v
Remove-Item Env:TEST_DATABASE_URL
```

The integration test covers health, homepage, typo search, product details, offers, duplicate registration, login, unauthenticated requests, quoting, price mismatch rollback, successful purchase, order/key retrieval and logout. It has not been executed here.

Manual browser checks:
1. Home sections load independently; hero arrows change banner.
2. Search `eldn ring` and `ring`; filter Steam and EUROPE; sort by price; clear filters. Check no-match state.
3. Open product, compare three offers and add two different sellers to cart; change quantity, remove item, refresh.
4. As guest, checkout redirects to login; register and return to checkout. Confirm fake payment; order history shows key and total.
5. Register a second account and confirm it cannot see the first account's orders. Logout and attempt direct checkout/order URLs.
6. Test insufficient inventory and changed price, including concurrent purchase attempts. Check that no key is assigned twice and failed transactions leave stock unchanged.
7. Verify mobile layout, keyboard controls, image URLs and browser console; compare against actual Eneba design. Replace logo with verified official asset before final review.
8. Repeat purchase flow against deployed URL, check direct route reload and CORS.
