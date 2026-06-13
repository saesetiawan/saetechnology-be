package website_setting

import (
	"time"

	"github.com/bagusaditiasetiawan/saetechnology-be/internal/domain/jsonvalue"

	"github.com/google/uuid"
)

type WebsiteSetting struct {
	ID                 uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	SiteName           string         `json:"site_name" gorm:"column:site_name"`
	Tagline            string         `json:"tagline" gorm:"column:tagline"`
	LogoURL            string         `json:"logo_url" gorm:"column:logo_url"`
	FaviconURL         string         `json:"favicon_url" gorm:"column:favicon_url"`
	PrimaryImageURL    string         `json:"primary_image_url" gorm:"column:primary_image_url"`
	SecondaryImageURL  string         `json:"secondary_image_url" gorm:"column:secondary_image_url"`
	BackgroundImageURL string         `json:"background_image_url" gorm:"column:background_image_url"`
	Email              string         `json:"email" gorm:"column:email"`
	Phone              string         `json:"phone" gorm:"column:phone"`
	Address            string         `json:"address" gorm:"column:address"`
	FacebookURL        string         `json:"facebook_url" gorm:"column:facebook_url"`
	InstagramURL       string         `json:"instagram_url" gorm:"column:instagram_url"`
	TiktokURL          string         `json:"tiktok_url" gorm:"column:tiktok_url"`
	PrimaryColor       string         `json:"primary_color" gorm:"column:primary_color"`
	SecondaryColor     string         `json:"secondary_color" gorm:"column:secondary_color"`
	AccentColor        string         `json:"accent_color" gorm:"column:accent_color"`
	BackgroundColor    string         `json:"background_color" gorm:"column:background_color"`
	SurfaceColor       string         `json:"surface_color" gorm:"column:surface_color"`
	TextColor          string         `json:"text_color" gorm:"column:text_color"`
	MutedTextColor     string         `json:"muted_text_color" gorm:"column:muted_text_color"`
	BorderColor        string         `json:"border_color" gorm:"column:border_color"`
	PrimaryContrastColor string       `json:"primary_contrast_color" gorm:"column:primary_contrast_color"`
	AccentContrastColor  string       `json:"accent_contrast_color" gorm:"column:accent_contrast_color"`
	SurfaceContrastColor string       `json:"surface_contrast_color" gorm:"column:surface_contrast_color"`
	SuccessColor         string       `json:"success_color" gorm:"column:success_color"`
	WarningColor         string       `json:"warning_color" gorm:"column:warning_color"`
	DangerColor          string       `json:"danger_color" gorm:"column:danger_color"`
	InfoColor            string       `json:"info_color" gorm:"column:info_color"`
	LabelColor           string       `json:"label_color" gorm:"column:label_color"`
	LabelBackgroundColor string       `json:"label_background_color" gorm:"column:label_background_color"`
	FontFamily           string       `json:"font_family" gorm:"column:font_family"`
	HeadingFontFamily    string       `json:"heading_font_family" gorm:"column:heading_font_family"`
	BorderRadius         string       `json:"border_radius" gorm:"column:border_radius"`
	ButtonRadius         string       `json:"button_radius" gorm:"column:button_radius"`
	ShadowStyle          string       `json:"shadow_style" gorm:"column:shadow_style"`
	Metadata           jsonvalue.JSON `json:"metadata" gorm:"column:metadata;type:jsonb"`
	CreatedAt          time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt          time.Time      `json:"updated_at" gorm:"column:updated_at"`
}

func (WebsiteSetting) TableName() string {
	return "website_settings"
}
