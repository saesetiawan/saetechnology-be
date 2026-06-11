package website_setting

type UpdateWebsiteSettingDto struct {
	SiteName           string `json:"site_name" validate:"required,max=120"`
	Tagline            string `json:"tagline"`
	LogoURL            string `json:"logo_url"`
	FaviconURL         string `json:"favicon_url"`
	PrimaryImageURL    string `json:"primary_image_url"`
	SecondaryImageURL  string `json:"secondary_image_url"`
	BackgroundImageURL string `json:"background_image_url"`
	Email              string `json:"email" validate:"omitempty,email,max=160"`
	Phone              string `json:"phone" validate:"max=50"`
	Address            string `json:"address"`
	FacebookURL        string `json:"facebook_url"`
	InstagramURL       string `json:"instagram_url"`
	TiktokURL          string `json:"tiktok_url"`
	PrimaryColor       string `json:"primary_color" validate:"max=20"`
	SecondaryColor     string `json:"secondary_color" validate:"max=20"`
	AccentColor        string `json:"accent_color" validate:"max=20"`
	BackgroundColor    string `json:"background_color" validate:"max=20"`
	SurfaceColor       string `json:"surface_color" validate:"max=20"`
	TextColor          string `json:"text_color" validate:"max=20"`
	MutedTextColor     string `json:"muted_text_color" validate:"max=20"`
	BorderColor        string `json:"border_color" validate:"max=20"`
	Metadata           string `json:"metadata"`
}
