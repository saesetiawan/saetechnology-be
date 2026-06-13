package product

import (
	"time"

	"github.com/bagusaditiasetiawan/saetechnology-be/internal/domain/jsonvalue"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	Slug        string         `json:"slug" gorm:"column:slug"`
	Name        string         `json:"name" gorm:"column:name"`
	Tagline     string         `json:"tagline" gorm:"column:tagline"`
	Summary     string         `json:"summary" gorm:"column:summary"`
	Description string         `json:"description" gorm:"column:description"`
	Category    string         `json:"category" gorm:"column:category"`
	Status      string         `json:"status" gorm:"column:status"`
	PriceLabel  string         `json:"price_label" gorm:"column:price_label"`
	PriceURL    string         `json:"price_url" gorm:"column:price_url"`
	ImageURL    string         `json:"image_url" gorm:"column:image_url"`
	DemoURL     string         `json:"demo_url" gorm:"column:demo_url"`
	CTALabel    string         `json:"cta_label" gorm:"column:cta_label"`
	CTAURL      string         `json:"cta_url" gorm:"column:cta_url"`
	SortOrder   int            `json:"sort_order" gorm:"column:sort_order"`
	IsFeatured  bool           `json:"is_featured" gorm:"column:is_featured"`
	Metadata    jsonvalue.JSON `json:"metadata" gorm:"column:metadata;type:jsonb"`
	CreatedAt   time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`
}

func (Product) TableName() string {
	return "products"
}
