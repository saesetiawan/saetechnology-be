package product

type CreateProductDto struct {
	Slug        string `json:"slug" validate:"required,max=140"`
	Name        string `json:"name" validate:"required,max=180"`
	Tagline     string `json:"tagline"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Category    string `json:"category" validate:"max=120"`
	Status      string `json:"status" validate:"required,oneof=draft published archived"`
	PriceLabel  string `json:"price_label" validate:"max=120"`
	PriceURL    string `json:"price_url"`
	ImageURL    string `json:"image_url"`
	DemoURL     string `json:"demo_url"`
	CTALabel    string `json:"cta_label" validate:"max=120"`
	CTAURL      string `json:"cta_url"`
	SortOrder   int    `json:"sort_order"`
	IsFeatured  bool   `json:"is_featured"`
	Metadata    string `json:"metadata"`
}

type UpdateProductDto struct {
	ID          string `json:"id"`
	Slug        string `json:"slug" validate:"required,max=140"`
	Name        string `json:"name" validate:"required,max=180"`
	Tagline     string `json:"tagline"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Category    string `json:"category" validate:"max=120"`
	Status      string `json:"status" validate:"required,oneof=draft published archived"`
	PriceLabel  string `json:"price_label" validate:"max=120"`
	PriceURL    string `json:"price_url"`
	ImageURL    string `json:"image_url"`
	DemoURL     string `json:"demo_url"`
	CTALabel    string `json:"cta_label" validate:"max=120"`
	CTAURL      string `json:"cta_url"`
	SortOrder   int    `json:"sort_order"`
	IsFeatured  bool   `json:"is_featured"`
	Metadata    string `json:"metadata"`
}

type ListProductQuery struct {
	Page       int
	Limit      int
	Search     string
	Status     string
	Category   string
	IsFeatured *bool
	OrderBy    string
	OrderType  string
	PublicOnly bool
}
