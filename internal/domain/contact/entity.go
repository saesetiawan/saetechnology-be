package contact

import (
	"time"

	"github.com/bagusaditiasetiawan/saetechnology-be/internal/domain/jsonvalue"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContactMessage struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	Name      string         `json:"name" gorm:"column:name"`
	Email     string         `json:"email" gorm:"column:email"`
	Phone     string         `json:"phone" gorm:"column:phone"`
	Company   string         `json:"company" gorm:"column:company"`
	Subject   string         `json:"subject" gorm:"column:subject"`
	Message   string         `json:"message" gorm:"column:message"`
	Source    string         `json:"source" gorm:"column:source"`
	Status    string         `json:"status" gorm:"column:status"`
	Metadata  jsonvalue.JSON `json:"metadata" gorm:"column:metadata;type:jsonb"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`
}

func (ContactMessage) TableName() string {
	return "contact_messages"
}
