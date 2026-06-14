package content

type CreateWebsiteContentDto struct {
	Key            string `json:"key" validate:"required,max=120"`
	Type           string `json:"type" validate:"required,oneof=hero banner section promo announcement carousel profile service service-grid service-detail use-case use-case-grid use-case-detail blog blog-grid blog-post process pricing faq testimonial case-study-grid"`
	Placement      string `json:"placement" validate:"required,max=80"`
	Title          string `json:"title" validate:"required,max=255"`
	Subtitle       string `json:"subtitle"`
	Body           string `json:"body"`
	ImageURL       string `json:"image_url"`
	LinkURL        string `json:"link_url"`
	LinkLabel      string `json:"link_label" validate:"max=120"`
	SortOrder      int    `json:"sort_order"`
	IsActive       bool   `json:"is_active"`
	Metadata       string `json:"metadata"`
	PublishStartAt string `json:"publish_start_at"`
	PublishEndAt   string `json:"publish_end_at"`
}

type UpdateWebsiteContentDto struct {
	ID             string `json:"id"`
	Key            string `json:"key" validate:"required,max=120"`
	Type           string `json:"type" validate:"required,oneof=hero banner section promo announcement carousel profile service service-grid service-detail use-case use-case-grid use-case-detail blog blog-grid blog-post process pricing faq testimonial case-study-grid"`
	Placement      string `json:"placement" validate:"required,max=80"`
	Title          string `json:"title" validate:"required,max=255"`
	Subtitle       string `json:"subtitle"`
	Body           string `json:"body"`
	ImageURL       string `json:"image_url"`
	LinkURL        string `json:"link_url"`
	LinkLabel      string `json:"link_label" validate:"max=120"`
	SortOrder      int    `json:"sort_order"`
	IsActive       bool   `json:"is_active"`
	Metadata       string `json:"metadata"`
	PublishStartAt string `json:"publish_start_at"`
	PublishEndAt   string `json:"publish_end_at"`
}

type ListWebsiteContentQuery struct {
	Page      int
	Limit     int
	Search    string
	SearchBy  string
	OrderBy   string
	OrderType string
	Type      string
	Placement string
	Active    *bool
}
