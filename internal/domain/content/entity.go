package content

import (
	"time"

	"go-platform-core/internal/domain/jsonvalue"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WebsiteContent struct {
	ID             uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	Key            string         `json:"key" gorm:"column:key"`
	Type           string         `json:"type" gorm:"column:type"`
	Placement      string         `json:"placement" gorm:"column:placement"`
	Title          string         `json:"title" gorm:"column:title"`
	Subtitle       string         `json:"subtitle" gorm:"column:subtitle"`
	Body           string         `json:"body" gorm:"column:body"`
	ImageURL       string         `json:"image_url" gorm:"column:image_url"`
	LinkURL        string         `json:"link_url" gorm:"column:link_url"`
	LinkLabel      string         `json:"link_label" gorm:"column:link_label"`
	SortOrder      int            `json:"sort_order" gorm:"column:sort_order"`
	IsActive       bool           `json:"is_active" gorm:"column:is_active"`
	Metadata       jsonvalue.JSON `json:"metadata" gorm:"column:metadata;type:jsonb"`
	PublishStartAt *time.Time     `json:"publish_start_at" gorm:"column:publish_start_at"`
	PublishEndAt   *time.Time     `json:"publish_end_at" gorm:"column:publish_end_at"`
	CreatedAt      time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`
}

func (WebsiteContent) TableName() string {
	return "website_contents"
}
