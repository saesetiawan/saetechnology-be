package website_setting

import (
	"time"

	"saetechnology-be/internal/domain/jsonvalue"

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
	Metadata           jsonvalue.JSON `json:"metadata" gorm:"column:metadata;type:jsonb"`
	CreatedAt          time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt          time.Time      `json:"updated_at" gorm:"column:updated_at"`
}

func (WebsiteSetting) TableName() string {
	return "website_settings"
}
