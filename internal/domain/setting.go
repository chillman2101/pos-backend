package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Setting struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Key       string         `gorm:"uniqueIndex;not null;size:100" json:"key"`
	Value     string         `gorm:"type:text" json:"value"`
	Category  string         `gorm:"size:50" json:"category"` // store, tax, receipt, system
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type SettingRepository interface {
	GetByKey(key string) (*Setting, error)
	GetByCategory(category string) ([]Setting, error)
	GetAll() ([]Setting, error)
	Upsert(setting *Setting) error
	BulkUpsert(settings []Setting) error
	Delete(key string) error
}
