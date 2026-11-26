package products

type Product struct {
	ID            int64  `json:"id"`
	Slug          string `json:"slug"`
	Name          string `json:"name"`
	Platform      string `json:"platform"`
	Region        string `json:"region"`
	Description   string `json:"description"`
	Image         string `json:"image"`
	OriginalCents int    `json:"originalCents"`
	PriceCents    int    `json:"priceCents"`
	Popularity    int    `json:"popularity"`
}
type Filter struct {
	Query, Platform, Region, Sort string
	Max                           int
	Limit                         int
}
