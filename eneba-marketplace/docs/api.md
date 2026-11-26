# REST API

Base URL: `http://localhost:8080/api`. JSON bodies and responses. Protected routes use `Authorization: Bearer <token>` returned by login/register. Monetary values are EUR cents.

| Method | Path | Input | Auth |
|---|---|---|---|
| GET | /health | none | no |
| GET | /products | q, platform, region, max (cents), sort | no |
| GET | /products/{slug} | slug | no |
| GET | /products/{slug}/offers | slug | no |
| GET | /homepage/banners | none | no |
| GET | /homepage/best-sellers | none | no |
| GET | /homepage/deals | none | no |
| POST | /auth/register | email, password | no |
| POST | /auth/login | email, password | no |
| GET | /auth/me | none | yes |
| POST | /auth/logout | none | yes |
| POST | /cart/quote | items | no |
| POST | /checkout | items, expectedTotal, payment | yes |
| GET | /orders | none | yes |

Sort: popularity (default), price_asc, price_desc, name, deals. Partial search uses LIKE; typo tolerance uses PostgreSQL word_similarity. Cart item: `{ "offerId": 1, "quantity": 1 }`. Quote returns lines and totalCents. Checkout accepts `payment: "demo"` and expectedTotal from the quote; returns orderId. Orders include demo keys belonging only to the authenticated user. Errors use `{ "error": "message" }`.
