package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"not null;size:255" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type CategoryRepository interface {
	Create(user *Category) error
	FindByID(id uuid.UUID) (*Category, error)
	FindByName(name string) (*Category, error)
	Update(category *Category) error
	Delete(id uuid.UUID) error
	FindAll(page, limit int) ([]Category, int64, error)
}
